import React from "react"

const ITEMS = [
  {
    id: "1",
    title: "アイテム 1",
    description: "説明テキストがここに入ります。",
  },
  {
    id: "2",
    title: "アイテム 2",
    description: "説明テキストがここに入ります。",
  },
  {
    id: "3",
    title: "アイテム 3",
    description: "説明テキストがここに入ります。",
  },
  {
    id: "4",
    title: "アイテム 4",
    description: "説明テキストがここに入ります。",
  },
  {
    id: "5",
    title: "アイテム 5",
    description: "説明テキストがここに入ります。",
  },
  {
    id: "6",
    title: "アイテム 6",
    description: "説明テキストがここに入ります。",
  },
]

export default function TopPage(): React.JSX.Element {
  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-2xl font-bold text-gray-100 mb-6">トップページ</h1>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          {ITEMS.map((item) => (
            <div
              key={item.id}
              className="bg-gray-800 rounded-xl shadow-sm border border-gray-700 p-5 hover:shadow-md transition-shadow"
            >
              <h2 className="text-base font-semibold text-gray-100 mb-2">
                {item.title}
              </h2>
              <p className="text-sm text-gray-400">{item.description}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
