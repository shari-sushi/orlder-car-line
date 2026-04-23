import { useParams } from "react-router-dom"
import type { ReactNode } from "react"

const pages = import.meta.glob("../../../docs/pages/*.md", {
  query: "?raw",
  import: "default",
  eager: true,
}) as Record<string, string>

type Token =
  | { type: "h1" | "h2" | "h3"; text: string }
  | { type: "ul_item"; text: string; indent: number }
  | { type: "table"; rows: string[][] }
  | { type: "img"; src: string; alt: string; height?: string }
  | { type: "p"; text: string }
  | { type: "blank" }

function parseInline(text: string): ReactNode {
  const parts: ReactNode[] = []
  const linkRe = /\[([^\]]+)\]\(([^)]+)\)|<(https?:\/\/[^>]+)>/g
  let last = 0
  let m: RegExpExecArray | null
  while ((m = linkRe.exec(text)) !== null) {
    if (m.index > last) parts.push(text.slice(last, m.index))
    const href = m[2] ?? m[3]
    const label = m[1] ?? href
    parts.push(
      <a
        key={m.index}
        href={href}
        target="_blank"
        rel="noopener noreferrer"
        className="text-blue-400 underline break-all"
      >
        {label}
      </a>,
    )
    last = m.index + m[0].length
  }
  if (last < text.length) parts.push(text.slice(last))
  return parts.length === 1 ? parts[0] : <>{parts}</>
}

function tokenize(md: string): Token[] {
  const tokens: Token[] = []
  const lines = md.split("\n")

  let i = 0
  while (i < lines.length) {
    const line = lines[i]

    if (line.trimStart().startsWith("|")) {
      const tableLines: string[] = []
      while (i < lines.length && lines[i].trimStart().startsWith("|")) {
        tableLines.push(lines[i])
        i++
      }
      const rows = tableLines
        .filter((l) => !/^\s*\|[\s|:-]+\|\s*$/.test(l))
        .map((l) =>
          l
            .replace(/^\s*\|/, "")
            .replace(/\|\s*$/, "")
            .split("|")
            .map((c) => c.trim()),
        )
      tokens.push({ type: "table", rows })
      continue
    }

    const h1 = line.match(/^# (.+)/)
    const h2 = line.match(/^## (.+)/)
    const h3 = line.match(/^### (.+)/)
    const ul = line.match(/^(\s*)[-*] (.+)/)
    const ulDot = line.match(/^(\s*)・(.+)/)
    const imgTag = line.match(/^<img\s([^>]*)>/)

    if (h1) tokens.push({ type: "h1", text: h1[1] })
    else if (h2) tokens.push({ type: "h2", text: h2[1] })
    else if (h3) tokens.push({ type: "h3", text: h3[1] })
    else if (ul)
      tokens.push({ type: "ul_item", text: ul[2], indent: ul[1].length })
    else if (ulDot)
      tokens.push({ type: "ul_item", text: ulDot[2], indent: ulDot[1].length })
    else if (imgTag) {
      const attrs = imgTag[1]
      tokens.push({
        type: "img",
        src: attrs.match(/src="([^"]+)"/)?.[1] ?? "",
        alt: attrs.match(/alt="([^"]*)"/)?.[1] ?? "",
        height: attrs.match(/height="([^"]*)"/)?.[1],
      })
    } else if (line.trim() === "") tokens.push({ type: "blank" })
    else tokens.push({ type: "p", text: line.trim() })

    i++
  }
  return tokens
}

function renderTokens(tokens: Token[]): ReactNode[] {
  const nodes: ReactNode[] = []
  let j = 0

  while (j < tokens.length) {
    const t = tokens[j]

    if (t.type === "h1") {
      nodes.push(
        <h1
          key={j}
          className="text-2xl font-bold text-gray-100 mt-2 mb-6 pb-2 border-b border-gray-700"
        >
          {parseInline(t.text)}
        </h1>,
      )
    } else if (t.type === "h2") {
      nodes.push(
        <h2
          key={j}
          className="text-xl font-bold text-gray-200 mt-8 mb-4 pb-1 border-b border-gray-700"
        >
          {parseInline(t.text)}
        </h2>,
      )
    } else if (t.type === "h3") {
      nodes.push(
        <h3 key={j} className="text-base font-semibold text-gray-300 mt-5 mb-2">
          {parseInline(t.text)}
        </h3>,
      )
    } else if (t.type === "ul_item") {
      const items: { text: string; indent: number }[] = []
      while (j < tokens.length && tokens[j].type === "ul_item") {
        const item = tokens[j] as {
          type: "ul_item"
          text: string
          indent: number
        }
        items.push({ text: item.text, indent: item.indent })
        j++
      }
      nodes.push(
        <ul key={`ul-${j}`} className="list-disc pl-5 my-2 space-y-0.5">
          {items.map((item, idx) => (
            <li
              key={idx}
              className="text-gray-300 text-sm"
              style={{ marginLeft: item.indent > 0 ? "1rem" : undefined }}
            >
              {parseInline(item.text)}
            </li>
          ))}
        </ul>,
      )
      continue
    } else if (t.type === "table") {
      const [header, ...body] = t.rows
      nodes.push(
        <div key={j} className="overflow-x-auto my-4">
          <table className="w-full border-collapse text-sm">
            <thead>
              <tr>
                {header.map((cell, ci) => (
                  <th
                    key={ci}
                    className="bg-gray-800 border border-gray-600 px-3 py-2 text-left font-semibold text-gray-300"
                  >
                    {parseInline(cell)}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {body.map((row, ri) => (
                <tr key={ri} className={ri % 2 === 0 ? "" : "bg-gray-800"}>
                  {row.map((cell, ci) => (
                    <td
                      key={ci}
                      className="border border-gray-600 px-3 py-1.5 text-gray-300"
                    >
                      {parseInline(cell)}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>,
      )
    } else if (t.type === "img") {
      nodes.push(
        <img
          key={j}
          src={`/${t.src}`}
          alt={t.alt}
          style={t.height ? { height: Number(t.height) } : undefined}
          className="my-2 max-w-full"
        />,
      )
    } else if (t.type === "p") {
      nodes.push(
        <p key={j} className="text-gray-300 text-sm leading-relaxed my-1">
          {parseInline(t.text)}
        </p>,
      )
    }

    j++
  }
  return nodes
}

export default function MdPage(): ReactNode {
  const { slug } = useParams<{ slug: string }>()

  const entry = Object.entries(pages).find(([path]) => {
    const filename = path.split("/").pop()?.replace(/\.md$/, "")
    return filename === slug
  })

  if (!entry) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-8 text-gray-400">
        ページが見つかりません: {slug}
      </div>
    )
  }

  const tokens = tokenize(entry[1])
  const nodes = renderTokens(tokens)

  return <div className="max-w-3xl mx-auto px-4 py-8">{nodes}</div>
}
