import { normalizeStockCode } from '@/utils/format'

export function moveWatchlistCodeBefore(codes: string[], draggedCode: string, targetCode: string): string[] {
  const dragged = normalizeStockCode(draggedCode)
  const target = normalizeStockCode(targetCode)
  if (!dragged || !target || dragged === target) return codes

  const normalized = codes.map(normalizeStockCode)
  const dragIndex = normalized.indexOf(dragged)
  const targetIndex = normalized.indexOf(target)
  if (dragIndex < 0 || targetIndex < 0) return codes

  const next = [...normalized]
  const [moved] = next.splice(dragIndex, 1)
  const nextTargetIndex = next.indexOf(target)
  next.splice(nextTargetIndex, 0, moved)
  return next
}
