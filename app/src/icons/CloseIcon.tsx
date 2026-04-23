interface CloseIconProps {
  width?: number
  height?: number
  className?: string
}

export default function ClosedIcon({
  className = "fill bg-zinc-600",
  height = 24,
  width = 24,
}: CloseIconProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      height={height}
      viewBox="0 -960 960 960"
      width={width}
      className={className}
    >
      <path d="m256-200-56-56 224-224-224-224 56-56 224 224 224-224 56 56-224 224 224 224-56 56-224-224-224 224Z" />
    </svg>
  )
}
