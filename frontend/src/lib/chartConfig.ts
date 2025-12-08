import { EChartsOption } from 'echarts'
import { ChartDataPoint, ChartSeries, Token } from './types'

// Color palette for different tokens
const TOKEN_COLORS: Record<string, string> = {
  wstETH: '#5470c6',
  ankrETH: '#91cc75',
  rETH: '#fac858',
  wBETH: '#ee6666',
  pufETH: '#73c0de',
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
    day: 'numeric',
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

  return { xAxisData, series }
}

export function createChartOption(
  xAxisData: string[],
  series: ChartSeries[],
  loading = false
): EChartsOption {
  return {
    title: {
      text: 'LST Token Price Comparison (vs ETH)',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
      },
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        if (!params || params.length === 0) return ''

        const date = params[0].name
        let content = `<strong>${date}</strong><br/>`

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
      top: 30,
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
      },
    },
    yAxis: {
      type: 'value',
      name: 'Price (ETH)',
      nameLocation: 'middle',
      nameGap: 50,
      min: 'dataMin',
      axisLabel: {
        formatter: (value: number) => formatPrice(value),
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
