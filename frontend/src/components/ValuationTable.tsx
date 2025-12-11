'use client'

import { useState, useMemo } from 'react'
import { ValuationData } from '@/lib/types'

interface ValuationTableProps {
  valuations: ValuationData[]
  loading?: boolean
}

type SortField = 'token_symbol' | 'price' | 'apr' | 'stability' | 'tvl' | 'remarks'
type SortDirection = 'asc' | 'desc'

export default function ValuationTable({ valuations, loading = false }: ValuationTableProps) {
  const [sortField, setSortField] = useState<SortField>('apr')
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc')

  const sortedValuations = useMemo(() => {
    return [...valuations].sort((a, b) => {
      const aValue = a[sortField]
      const bValue = b[sortField]

      // Handle string sorting
      if (typeof aValue === 'string' && typeof bValue === 'string') {
        const aStr = aValue.toLowerCase()
        const bStr = bValue.toLowerCase()
        if (aStr < bStr) return sortDirection === 'asc' ? -1 : 1
        if (aStr > bStr) return sortDirection === 'asc' ? 1 : -1
        return 0
      }

      // Handle number sorting
      if (typeof aValue === 'number' && typeof bValue === 'number') {
        if (aValue < bValue) return sortDirection === 'asc' ? -1 : 1
        if (aValue > bValue) return sortDirection === 'asc' ? 1 : -1
        return 0
      }

      return 0
    })
  }, [valuations, sortField, sortDirection])

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
    } else {
      setSortField(field)
      setSortDirection('desc') // Default to descending for new field
    }
  }

  const formatTVL = (tvl: number) => {
    if (tvl >= 1e9) return `${(tvl / 1e9).toFixed(1)}B`
    if (tvl >= 1e6) return `${(tvl / 1e6).toFixed(1)}M`
    if (tvl >= 1e3) return `${(tvl / 1e3).toFixed(1)}K`
    return `${tvl.toFixed(0)}`
  }

  const getRemarksColor = (remarks: string) => {
    switch (remarks) {
      case 'Very Undervalued':
      case 'Undervalued':
        return 'text-green-600 dark:text-green-400'
      case 'Very Overvalued':
      case 'Overvalued':
        return 'text-red-600 dark:text-red-400'
      case 'Fair Value':
        return 'text-blue-600 dark:text-blue-400'
      default:
        return 'text-gray-600 dark:text-gray-400'
    }
  }

  const SortIcon = ({ field }: { field: SortField }) => {
    if (sortField !== field) return null
    return (
      <span className="ml-1">
        {sortDirection === 'asc' ? '↑' : '↓'}
      </span>
    )
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-48">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 mx-auto mb-2" style={{ borderBottomColor: 'rgb(59, 130, 246)' }}></div>
          <p style={{ color: 'rgb(var(--text-secondary))' }}>Loading valuation data...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead style={{ backgroundColor: 'rgb(var(--bg-secondary))' }}>
          <tr>
            <th
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800"
              onClick={() => handleSort('token_symbol')}
              style={{ color: 'rgb(var(--text-secondary))' }}
            >
              Token <SortIcon field="token_symbol" />
            </th>
            <th
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800"
              onClick={() => handleSort('price')}
              style={{ color: 'rgb(var(--text-secondary))' }}
            >
              Price <SortIcon field="price" />
            </th>
            <th
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800"
              onClick={() => handleSort('apr')}
              style={{ color: 'rgb(var(--text-secondary))' }}
            >
              APR <SortIcon field="apr" />
            </th>
            <th
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800"
              onClick={() => handleSort('stability')}
              style={{ color: 'rgb(var(--text-secondary))' }}
            >
              Stability <SortIcon field="stability" />
            </th>
            <th
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800"
              onClick={() => handleSort('tvl')}
              style={{ color: 'rgb(var(--text-secondary))' }}
            >
              TVL <SortIcon field="tvl" />
            </th>
            <th
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800"
              onClick={() => handleSort('remarks')}
              style={{ color: 'rgb(var(--text-secondary))' }}
            >
              Remarks (vs price right now) <SortIcon field="remarks" />
            </th>
          </tr>
        </thead>
        <tbody style={{ backgroundColor: 'rgb(var(--bg-primary))' }} className="divide-y divide-gray-200 dark:divide-gray-700">
          {sortedValuations.map((valuation) => (
            <tr key={valuation.token_symbol} className="hover:bg-gray-50 dark:hover:bg-gray-800/50">
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium" style={{ color: 'rgb(var(--text-primary))' }}>
                {valuation.token_symbol}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm" style={{ color: 'rgb(var(--text-primary))' }}>
                {valuation.price.toFixed(4)} ETH
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm" style={{ color: 'rgb(var(--text-primary))' }}>
                {(valuation.apr * 100).toFixed(2)}%
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm" style={{ color: 'rgb(var(--text-primary))' }}>
                {Math.max(0, Math.min(100, valuation.stability * 3000)).toFixed(1)}%
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm" style={{ color: 'rgb(var(--text-primary))' }}>
                {formatTVL(valuation.tvl)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm">
                <span className={`font-medium ${getRemarksColor(valuation.remarks)}`}>
                  {valuation.remarks}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
