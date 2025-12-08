'use client'

import React, { useMemo } from 'react'
import ReactECharts from 'echarts-for-react'
import { Token } from '@/lib/types'
import { prepareChartData, createChartOption } from '@/lib/chartConfig'

interface ChartProps {
  tokens: Token[]
  histories: Record<string, { price_history: Array<{ timestamp: number; price: number }> }>
  loading?: boolean
  className?: string
}

export default function Chart({ tokens, histories, loading = false, className = '' }: ChartProps) {
  const chartOption = useMemo(() => {
    if (tokens.length === 0) {
      return createChartOption([], [], loading)
    }

    const { xAxisData, series } = prepareChartData(tokens, histories)
    return createChartOption(xAxisData, series, loading)
  }, [tokens, histories, loading])

  const chartStyle = {
    height: '500px',
    width: '100%',
  }

  return (
    <div className={`w-full ${className}`}>
      <ReactECharts
        option={chartOption}
        style={chartStyle}
        opts={{ renderer: 'canvas' }}
        className="w-full"
      />
    </div>
  )
}
