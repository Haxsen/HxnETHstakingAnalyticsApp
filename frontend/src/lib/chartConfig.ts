import { EChartsOption } from 'echarts'
import { ChartDataPoint, ChartSeries, Token } from './types'

// Color palette for different tokens
const TOKEN_COLORS: Record<string, string> = {
  wstETH: '#5470c6',
  ankrETH: '#91cc75',
  rETH: '#fac858',
  wBETH: '#ee6666',
  pufETH: '#73c0de',
  LSETH: '#ff7c7c',
  RSETH: '#ffb347',
  METH: '#87ceeb',
  CBETH: '#dda0dd',
  TETH: '#98fb98',
  SFRXETH: '#f0e68c',
  CDCETH: '#ffa07a',
  UNIETH: '#20b2aa',
  ETH: '#3ba272', // Reference ETH color
}

// Default color for unknown tokens
const DEFAULT_COLOR = '#cccccc'

export function getTokenColor(symbol: string): string {
  return TOKEN_COLORS[symbol] || DEFAULT_COLOR
}

export function formatTimestamp(timestamp: number): string {
  return new Date(timestamp).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
  })
}

export function formatPrice(price: number): string {
  return price.toFixed(6) // ETH prices are typically small decimals
}

export function prepareChartData(
  tokens: Token[],
  histories: Record<string, { price_history: Array<{ timestamp: number; price: number }> }>
): {
  xAxisData: string[]
  timestamps: number[]
  series: ChartSeries[]
} {
  // Collect all unique timestamps and sort them
  const allTimestamps = new Set<number>()

  Object.values(histories).forEach(history => {
    history.price_history.forEach(point => {
      allTimestamps.add(point.timestamp)
    })
  })

  const sortedTimestamps = Array.from(allTimestamps).sort((a, b) => a - b)
  const xAxisData = sortedTimestamps.map(formatTimestamp)

  // Create series for each token
  const series: ChartSeries[] = []

  tokens.forEach(token => {
    const history = histories[token.symbol]
    if (!history) return

    // Create a map of timestamp to price for quick lookup
    const priceMap = new Map<number, number>()
    history.price_history.forEach(point => {
      priceMap.set(point.timestamp, point.price)
    })

    // Create data array aligned with xAxisData
    const data = sortedTimestamps.map(timestamp => {
      return priceMap.get(timestamp) || null
    })

    series.push({
      name: token.symbol,
      type: 'line',
      data: data as number[], // ECharts handles null values
      smooth: false,
      symbol: 'none',
      lineStyle: {
        width: 2,
      },
      itemStyle: {
        color: getTokenColor(token.symbol),
      },
    })
  })

  // Add ETH reference line if available
  if (histories.ETH) {
    const ethHistory = histories.ETH
    const ethPriceMap = new Map<number, number>()
    ethHistory.price_history.forEach(point => {
      ethPriceMap.set(point.timestamp, point.price)
    })

    const ethData = sortedTimestamps.map(timestamp => {
      return ethPriceMap.get(timestamp) || null
    })

    series.push({
      name: 'ETH',
      type: 'line',
      data: ethData as number[],
      smooth: false,
      symbol: 'none',
      lineStyle: {
        width: 2,
        type: 'dashed',
      },
      itemStyle: {
        color: getTokenColor('ETH'),
      },
    })
  }

  return { xAxisData, timestamps: sortedTimestamps, series }
}

export function createChartOption(
  xAxisData: string[],
  timestamps: number[],
  series: ChartSeries[],
  loading = false,
  theme: 'light' | 'dark' = 'light'
): EChartsOption {
  const textColor = theme === 'dark' ? '#f8fafc' : '#0f172a' // slate-50 : slate-900
  const axisColor = theme === 'dark' ? '#64748b' : '#64748b' // slate-500 for both (good contrast)

  return {
    title: {
      text: 'LST Token Price Comparison (vs ETH)',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: textColor,
      },
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any, ticket: any, callback: any) => {
        if (!params || params.length === 0) return ''

        // Get the data index to find the original timestamp
        const dataIndex = params[0].dataIndex
        const originalTimestamp = timestamps[dataIndex]

        // Format the full date for tooltip
        const fullDate = new Date(originalTimestamp).toLocaleDateString('en-US', {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
        })

        let content = `<strong>${fullDate}</strong><br/>`

        params.forEach((param: any) => {
          if (param.value !== null && param.value !== undefined) {
            content += `${param.seriesName}: ${formatPrice(param.value)} ETH<br/>`
          }
        })

        return content
      },
    },
    legend: {
      data: series.map(s => s.name),
      top: 50,
      textStyle: {
        color: textColor,
      },
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '15%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: xAxisData,
      axisLabel: {
        rotate: 45,
        interval: 'auto',
        color: axisColor,
      },
      nameTextStyle: {
        color: textColor,
      },
    },
    yAxis: {
      type: 'value',
      name: 'Price (ETH)',
      nameLocation: 'middle',
      nameGap: 70,
      min: 'dataMin',
      axisLabel: {
        formatter: (value: number) => formatPrice(value),
        color: axisColor,
      },
      nameTextStyle: {
        color: textColor,
      },
    },
    series,
    loading: loading,
    dataZoom: [
      {
        type: 'inside',
        start: 0,
        end: 100,
      },
      {
        start: 0,
        end: 100,
      },
    ],
    responsive: true,
  }
}
