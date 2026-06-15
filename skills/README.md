# 本地 Skills 使用说明

本目录保存策衡项目的本地 Claude Code skills，用于沉淀项目内可复用的工作流程。

## 目录约定

- 每个 skill 使用独立子目录，例如 `skills/frontend-design/`。
- skill 主说明写在子目录的 `SKILL.md` 中。
- `SKILL.md` 使用 frontmatter 声明 `name` 和 `description`。
- 项目内新增或调整 skill 时，应同步更新本说明或对应 skill 的 `SKILL.md`。

## 强制约束

涉及前端页面、组件、样式、交互、状态管理或前端数据展示的任务，写代码前必须先使用 `skills/frontend-design`。

使用该 skill 后，需要先产出“前端设计结论”，确认页面目标、范围、信息层级、数据依赖、交互、状态和验收方式，再进入实现。
