---
paths: ["**/*.md"]
---

# Markdownファイルのルール

## 基本

- すべてのMarkdownファイルは `markdownlint` のルールに準拠すること
- Markdownファイルを編集する際は、必ず `markdownlint` の警告を出さないように記述すること

## よく起こるミス

- MD022/blanks-around-headings
- MD031/blanks-around-fences
- MD036/no-emphasis-as-heading: 見出しの代わりに太字・斜体を使わない。見出しには `#` を使う
- MD040/fenced-code-language: Fenced code blocks should have a language
- MD047/single-trailing-newline: Files should end with a single newline character
- MD060/table-column-style: Table column style [Table pipe is missing space]

## 水平線（`---`）

明示的な指示がない限り `---` による水平線は使わない。セクションの区切りには見出し（`##` など）を使う。
