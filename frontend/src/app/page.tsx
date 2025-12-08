'use client'

import { useState, useEffect } from 'react'
import Chart from '@/components/Chart'
import ThemeToggle from '@/components/ThemeToggle'
import { Token, LoadingState, ErrorState } from '@/lib/types'
import { fetchTokens, fetchTokenHistory, getErrorMessage } from '@/lib/api'

export default function Home() {
  const [tokens, setTokens] = useState<Token[]>([])
  const [histories, setHistories] = useState<Record<string, any>>({})
  const [loading, setLoading] = useState<LoadingState>({
    tokens: true,
    history: false,
  })
  const [errors, setErrors] = useState<ErrorState>({})

  // Fetch tokens on mount
  useEffect(() => {
    const loadTokens = async () => {
      try {
        setLoading(prev => ({ ...prev, tokens: true }))
        setErrors(prev => ({ ...prev, tokens: undefined }))

        const response = await fetchTokens()
        setTokens(response.tokens)
      } catch (error) {
        setErrors(prev => ({ ...prev, tokens: getErrorMessage(error) }))
      } finally {
        setLoading(prev => ({ ...prev, tokens: false }))
      }
    }

    loadTokens()
  }, [])

  // Fetch price histories when tokens are loaded
  useEffect(() => {
    if (tokens.length === 0) return

    const loadHistories = async () => {
      setLoading(prev => ({ ...prev, history: true }))
      setErrors(prev => ({ ...prev, history: undefined }))

      const newHistories: Record<string, any> = {}

      try {
        // Fetch history for each token
        const promises = tokens.map(async (token) => {
          try {
            const response = await fetchTokenHistory(token.symbol)
            newHistories[token.symbol] = response
          } catch (error) {
            console.error(`Failed to fetch history for ${token.symbol}:`, error)
            // Continue with other tokens even if one fails
          }
        })

        await Promise.all(promises)
        setHistories(newHistories)
      } catch (error) {
        setErrors(prev => ({ ...prev, history: getErrorMessage(error) }))
      } finally {
        setLoading(prev => ({ ...prev, history: false }))
      }
    }

    loadHistories()
  }, [tokens])

  const isLoading = loading.tokens || loading.history
  const hasError = errors.tokens || errors.history

  return (
    <main className="min-h-screen py-8 transition-colors duration-200" style={{ backgroundColor: 'rgb(var(--bg-primary))' }}>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8 flex items-start justify-between">
          <div>
            <h1 className="text-3xl font-bold mb-2" style={{ color: 'rgb(var(--text-primary))' }}>
              ETH Staking Analytics Dashboard
            </h1>
            <p style={{ color: 'rgb(var(--text-secondary))' }}>
              Compare Liquid Staking Token performance against ETH over the past year
            </p>
          </div>
          <ThemeToggle />
        </div>

        {/* Error Display */}
        {hasError && (
          <div className="mb-6 p-4 rounded-md border" style={{
            backgroundColor: 'rgba(239, 68, 68, 0.1)',
            borderColor: 'rgb(220, 38, 38)',
          }}>
            <div className="flex">
              <div className="ml-3">
                <h3 className="text-sm font-medium" style={{ color: 'rgb(185, 28, 28)' }}>
                  Error loading data
                </h3>
                <div className="mt-2 text-sm" style={{ color: 'rgb(153, 27, 27)' }}>
                  {errors.tokens && <p>Tokens: {errors.tokens}</p>}
                  {errors.history && <p>Price data: {errors.history}</p>}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Chart Container */}
        <div className="rounded-lg shadow-sm border p-6 transition-colors duration-200" style={{
          backgroundColor: 'rgb(var(--bg-secondary))',
          borderColor: 'rgb(var(--border))'
        }}>
          {isLoading ? (
            <div className="flex items-center justify-center h-96">
              <div className="text-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 mx-auto mb-4" style={{ borderBottomColor: 'rgb(59, 130, 246)' }}></div>
                <p style={{ color: 'rgb(var(--text-secondary))' }}>Loading price data...</p>
              </div>
            </div>
          ) : tokens.length > 0 ? (
            <Chart
              tokens={tokens}
              histories={histories}
              loading={false}
            />
          ) : (
            <div className="flex items-center justify-center h-96">
              <div className="text-center">
                <p style={{ color: 'rgb(var(--text-secondary))' }}>No token data available</p>
              </div>
            </div>
          )}
        </div>

        {/* Token List */}
        {tokens.length > 0 && (
          <div className="mt-8">
            <h2 className="text-xl font-semibold mb-4" style={{ color: 'rgb(var(--text-primary))' }}>
              Tracked Tokens ({tokens.length})
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {tokens.map((token) => (
                <div
                  key={token.symbol}
                  className="rounded-lg shadow-sm border p-4 transition-colors duration-200"
                  style={{
                    backgroundColor: 'rgb(var(--bg-secondary))',
                    borderColor: 'rgb(var(--border))'
                  }}
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <h3 className="font-medium" style={{ color: 'rgb(var(--text-primary))' }}>{token.symbol}</h3>
                      <p className="text-sm" style={{ color: 'rgb(var(--text-secondary))' }}>{token.name}</p>
                    </div>
                    <div className="text-right">
                      <span className="text-xs uppercase tracking-wide" style={{ color: 'rgb(var(--text-secondary))' }}>
                        {token.blockchain}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </main>
  )
}
