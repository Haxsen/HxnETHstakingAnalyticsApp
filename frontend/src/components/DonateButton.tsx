'use client'

import { useState } from 'react'
import DonateModal from './DonateModal'

function DonateButton() {
  const [isModalOpen, setIsModalOpen] = useState(false)

  return (
    <>
      <button
        onClick={() => setIsModalOpen(true)}
        className="px-3 py-1.5 text-sm rounded transition-colors"
        style={{
          backgroundColor: 'rgb(var(--bg-secondary))',
          border: '1px solid rgb(var(--border))',
          color: 'rgb(var(--text-primary))',
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.backgroundColor = 'rgb(var(--border))'
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.backgroundColor = 'rgb(var(--bg-secondary))'
        }}
      >
        Donate
      </button>

      <DonateModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      />
    </>
  )
}

export default DonateButton
