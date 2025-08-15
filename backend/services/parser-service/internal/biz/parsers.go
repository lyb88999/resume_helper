package biz

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/unidoc/unioffice/document"
)

// TextParser 文本解析器
type TextParser struct{}

func NewTextParser() *TextParser {
	return &TextParser{}
}

func (p *TextParser) SupportedTypes() []string {
	return []string{"txt"}
}

func (p *TextParser) Parse(ctx context.Context, filePath string, options *ParseOptions) (*ParsedContent, error) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	text := string(content)
	if len(strings.TrimSpace(text)) == 0 {
		return nil, ErrEmptyContent
	}

	// 解析结构化内容
	parsedContent := &ParsedContent{
		RawText: text,
	}

	// 提取个人信息
	parsedContent.PersonalInfo = p.extractPersonalInfo(text)

	// 提取教育背景
	parsedContent.Education = p.extractEducation(text)

	// 提取工作经历
	parsedContent.Experience = p.extractExperience(text)

	// 提取技能
	parsedContent.Skills = p.extractSkills(text)

	// 提取项目经历
	parsedContent.Projects = p.extractProjects(text)

	return parsedContent, nil
}

// extractPersonalInfo 提取个人信息
func (p *TextParser) extractPersonalInfo(text string) *PersonalInfo {
	info := &PersonalInfo{}

	// 提取姓名 (通常在文档开头)
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 && i < 5 { // 前5行中寻找姓名
			// 简单的姓名匹配 (中文或英文)
			if matched, _ := regexp.MatchString(`^[\u4e00-\u9fa5a-zA-Z\s]{2,20}$`, line); matched {
				info.Name = line
				break
			}
		}
	}

	// 提取邮箱
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	if emails := emailRegex.FindAllString(text, -1); len(emails) > 0 {
		info.Email = emails[0]
	}

	// 提取电话号码
	phoneRegex := regexp.MustCompile(`1[3-9]\d{9}|(\+86)?[\s-]?1[3-9]\d{9}`)
	if phones := phoneRegex.FindAllString(text, -1); len(phones) > 0 {
		info.Phone = phones[0]
	}

	return info
}

// extractEducation 提取教育背景
func (p *TextParser) extractEducation(text string) []*Education {
	var educations []*Education

	// 寻找教育相关关键词
	eduKeywords := []string{"教育背景", "教育经历", "学历", "毕业", "大学", "学院", "专业"}
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		for _, keyword := range eduKeywords {
			if strings.Contains(line, keyword) {
				// 在附近几行寻找教育信息
				edu := p.parseEducationSection(lines, i, i+5)
				if edu != nil {
					educations = append(educations, edu)
				}
				break
			}
		}
	}

	return educations
}

func (p *TextParser) parseEducationSection(lines []string, start, end int) *Education {
	if end >= len(lines) {
		end = len(lines) - 1
	}

	edu := &Education{}
	section := strings.Join(lines[start:end+1], " ")

	// 提取学校名称
	schoolRegex := regexp.MustCompile(`([\u4e00-\u9fa5]{2,10}(大学|学院|学校))`)
	if schools := schoolRegex.FindAllString(section, -1); len(schools) > 0 {
		edu.School = schools[0]
	}

	// 提取专业
	majorRegex := regexp.MustCompile(`专业[:：]?\s*([\u4e00-\u9fa5a-zA-Z\s]{2,20})`)
	if majors := majorRegex.FindStringSubmatch(section); len(majors) > 1 {
		edu.Major = strings.TrimSpace(majors[1])
	}

	// 提取学位
	degreeKeywords := []string{"本科", "学士", "硕士", "博士", "专科"}
	for _, degree := range degreeKeywords {
		if strings.Contains(section, degree) {
			edu.Degree = degree
			break
		}
	}

	if edu.School != "" || edu.Major != "" {
		return edu
	}
	return nil
}

// extractExperience 提取工作经历
func (p *TextParser) extractExperience(text string) []*Experience {
	var experiences []*Experience

	// 寻找工作经历相关关键词
	expKeywords := []string{"工作经历", "工作经验", "职业经历", "任职", "工作"}
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		for _, keyword := range expKeywords {
			if strings.Contains(line, keyword) {
				// 在附近几行寻找工作信息
				exp := p.parseExperienceSection(lines, i, i+8)
				if exp != nil {
					experiences = append(experiences, exp)
				}
				break
			}
		}
	}

	return experiences
}

func (p *TextParser) parseExperienceSection(lines []string, start, end int) *Experience {
	if end >= len(lines) {
		end = len(lines) - 1
	}

	exp := &Experience{}
	section := strings.Join(lines[start:end+1], " ")

	// 提取公司名称
	companyRegex := regexp.MustCompile(`([\u4e00-\u9fa5a-zA-Z\s]{2,20}(公司|集团|科技|有限责任公司|股份有限公司))`)
	if companies := companyRegex.FindAllString(section, -1); len(companies) > 0 {
		exp.Company = companies[0]
	}

	// 提取职位
	positionKeywords := []string{"工程师", "经理", "主管", "总监", "专员", "助理", "开发", "设计师"}
	for _, pos := range positionKeywords {
		if strings.Contains(section, pos) {
			exp.Position = pos
			break
		}
	}

	if exp.Company != "" || exp.Position != "" {
		return exp
	}
	return nil
}

// extractSkills 提取技能
func (p *TextParser) extractSkills(text string) *Skills {
	skills := &Skills{
		Categories: []*SkillCategory{},
	}

	// 技术技能关键词
	techSkills := []string{
		"Java", "Python", "JavaScript", "Go", "C++", "C#", "PHP", "Ruby",
		"React", "Vue", "Angular", "Spring", "Django", "Flask",
		"MySQL", "PostgreSQL", "MongoDB", "Redis",
		"Docker", "Kubernetes", "Git", "Linux",
	}

	var foundSkills []*SkillItem
	for _, skill := range techSkills {
		if strings.Contains(text, skill) {
			foundSkills = append(foundSkills, &SkillItem{
				Name:  skill,
				Level: "熟练",
			})
		}
	}

	if len(foundSkills) > 0 {
		skills.Categories = append(skills.Categories, &SkillCategory{
			Category: "技术技能",
			Skills:   foundSkills,
		})
	}

	return skills
}

// extractProjects 提取项目经历
func (p *TextParser) extractProjects(text string) []*Project {
	var projects []*Project

	// 寻找项目相关关键词
	projectKeywords := []string{"项目经历", "项目经验", "参与项目", "负责项目"}
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		for _, keyword := range projectKeywords {
			if strings.Contains(line, keyword) {
				// 在附近几行寻找项目信息
				project := p.parseProjectSection(lines, i, i+6)
				if project != nil {
					projects = append(projects, project)
				}
				break
			}
		}
	}

	return projects
}

func (p *TextParser) parseProjectSection(lines []string, start, end int) *Project {
	if end >= len(lines) {
		end = len(lines) - 1
	}

	project := &Project{}
	section := strings.Join(lines[start:end+1], " ")

	// 提取项目名称 (通常在引号中或特定格式)
	nameRegex := regexp.MustCompile(`[""]([^"""]{3,30})[""]|项目[:：]\s*([\u4e00-\u9fa5a-zA-Z\s]{3,30})`)
	if names := nameRegex.FindStringSubmatch(section); len(names) > 1 {
		if names[1] != "" {
			project.Name = names[1]
		} else if names[2] != "" {
			project.Name = strings.TrimSpace(names[2])
		}
	}

	// 提取角色
	roleKeywords := []string{"负责人", "开发者", "架构师", "项目经理", "团队leader"}
	for _, role := range roleKeywords {
		if strings.Contains(section, role) {
			project.Role = role
			break
		}
	}

	if project.Name != "" {
		return project
	}
	return nil
}

// MarkdownParser Markdown解析器
type MarkdownParser struct {
	textParser *TextParser
}

func NewMarkdownParser() *MarkdownParser {
	return &MarkdownParser{
		textParser: NewTextParser(),
	}
}

func (p *MarkdownParser) SupportedTypes() []string {
	return []string{"md", "markdown"}
}

func (p *MarkdownParser) Parse(ctx context.Context, filePath string, options *ParseOptions) (*ParsedContent, error) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	text := string(content)
	if len(strings.TrimSpace(text)) == 0 {
		return nil, ErrEmptyContent
	}

	// 移除Markdown标记
	cleanedText := p.cleanMarkdown(text)

	// 使用文本解析器处理
	parsedContent, err := p.textParser.Parse(ctx, filePath, options)
	if err != nil {
		return nil, err
	}

	// 更新原始文本为清理后的内容
	parsedContent.RawText = cleanedText

	return parsedContent, nil
}

func (p *MarkdownParser) cleanMarkdown(text string) string {
	// 移除Markdown标记
	text = regexp.MustCompile(`#{1,6}\s*`).ReplaceAllString(text, "")          // 标题
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "$1")    // 粗体
	text = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(text, "$1")        // 斜体
	text = regexp.MustCompile(`\[(.*?)\]\(.*?\)`).ReplaceAllString(text, "$1") // 链接
	text = regexp.MustCompile("```[\\s\\S]*?```").ReplaceAllString(text, "")   // 代码块
	text = regexp.MustCompile("`(.*?)`").ReplaceAllString(text, "$1")          // 行内代码
	text = regexp.MustCompile(`(?m)^\s*[-*+]\s+`).ReplaceAllString(text, "")   // 列表标记

	return text
}

// PDFParser PDF解析器
type PDFParser struct {
	textParser *TextParser
}

func NewPDFParser() *PDFParser {
	return &PDFParser{
		textParser: NewTextParser(),
	}
}

func (p *PDFParser) SupportedTypes() []string {
	return []string{"pdf"}
}

func (p *PDFParser) Parse(ctx context.Context, filePath string, options *ParseOptions) (*ParsedContent, error) {
	// 使用github.com/ledongthuc/pdf解析PDF
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer file.Close()

	var textContent strings.Builder
	numPages := reader.NumPage()

	// 逐页提取文本
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			// 如果某页解析失败，继续处理其他页面
			continue
		}

		textContent.WriteString(text)
		textContent.WriteString("\n")
	}

	extractedText := textContent.String()
	if len(strings.TrimSpace(extractedText)) == 0 {
		return nil, ErrEmptyContent
	}

	// 使用文本解析器处理提取的文本
	tempFile := filePath + ".tmp.txt"
	err = ioutil.WriteFile(tempFile, []byte(extractedText), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(tempFile) // 清理临时文件

	parsedContent, err := p.textParser.Parse(ctx, tempFile, options)
	if err != nil {
		return nil, err
	}

	// 更新元数据
	if parsedContent.Metadata == nil {
		parsedContent.Metadata = &ParseMetadata{}
	}
	parsedContent.Metadata.PageCount = int32(numPages)
	parsedContent.Metadata.ParserVersion = "PDF-1.0.0"
	parsedContent.RawText = extractedText

	return parsedContent, nil
}

// DocxParser Word文档解析器
type DocxParser struct {
	textParser *TextParser
}

func NewDocxParser() *DocxParser {
	return &DocxParser{
		textParser: NewTextParser(),
	}
}

func (p *DocxParser) SupportedTypes() []string {
	return []string{"docx"}
}

func (p *DocxParser) Parse(ctx context.Context, filePath string, options *ParseOptions) (*ParsedContent, error) {
	// 只支持.docx格式，不支持老的.doc格式
	if !strings.HasSuffix(strings.ToLower(filePath), ".docx") {
		return nil, fmt.Errorf("unsupported format: only .docx files are supported")
	}

	// 使用github.com/unidoc/unioffice解析Word文档
	doc, err := document.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Word document: %w", err)
	}
	defer doc.Close()

	var textContent strings.Builder

	// 提取所有段落的文本
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			textContent.WriteString(run.Text())
		}
		textContent.WriteString("\n")
	}

	// 提取表格中的文本
	for _, table := range doc.Tables() {
		for _, row := range table.Rows() {
			for _, cell := range row.Cells() {
				for _, para := range cell.Paragraphs() {
					for _, run := range para.Runs() {
						textContent.WriteString(run.Text())
						textContent.WriteString(" ")
					}
				}
				textContent.WriteString("\t")
			}
			textContent.WriteString("\n")
		}
	}

	extractedText := textContent.String()
	if len(strings.TrimSpace(extractedText)) == 0 {
		return nil, ErrEmptyContent
	}

	// 使用文本解析器处理提取的文本
	tempFile := filePath + ".tmp.txt"
	err = ioutil.WriteFile(tempFile, []byte(extractedText), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(tempFile) // 清理临时文件

	parsedContent, err := p.textParser.Parse(ctx, tempFile, options)
	if err != nil {
		return nil, err
	}

	// 更新元数据
	if parsedContent.Metadata == nil {
		parsedContent.Metadata = &ParseMetadata{}
	}
	parsedContent.Metadata.ParserVersion = "DOCX-1.0.0"
	parsedContent.RawText = extractedText

	// 计算大概页数（基于字符数估算）
	textLength := len(extractedText)
	estimatedPages := (textLength / 1500) + 1 // 大概每页1500字符
	if estimatedPages < 1 {
		estimatedPages = 1
	}
	parsedContent.Metadata.PageCount = int32(estimatedPages)

	return parsedContent, nil
}
