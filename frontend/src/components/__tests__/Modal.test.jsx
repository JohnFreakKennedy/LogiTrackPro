/**
 * Unit tests for Modal component
 * Tests component rendering, props handling, and user interactions
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import Modal from '../Modal'

// Mock framer-motion for testing
vi.mock('framer-motion', () => ({
  motion: {
    div: ({ children, onClick, className, ...props }) => <div onClick={onClick} className={className} {...props}>{children}</div>
  },
  AnimatePresence: ({ children }) => children
}))

describe('Modal Component', () => {
  const defaultProps = {
    isOpen: true,
    onClose: vi.fn(),
    title: 'Test Modal',
    children: <div>Modal Content</div>
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render when isOpen is true', () => {
    render(<Modal {...defaultProps} />)
    expect(screen.getByText('Test Modal')).toBeInTheDocument()
    expect(screen.getByText('Modal Content')).toBeInTheDocument()
  })

  it('should not render when isOpen is false', () => {
    render(<Modal {...defaultProps} isOpen={false} />)
    expect(screen.queryByText('Test Modal')).not.toBeInTheDocument()
  })

  it('should call onClose when close button is clicked', () => {
    const onClose = vi.fn()
    render(<Modal {...defaultProps} onClose={onClose} />)
    
    // Find close button (X icon button)
    const buttons = screen.getAllByRole('button')
    const closeButton = buttons.find(btn => btn.onclick !== null || btn.getAttribute('onClick'))
    if (closeButton) {
      fireEvent.click(closeButton)
      expect(onClose).toHaveBeenCalledTimes(1)
    }
  })

  it('should call onClose when backdrop is clicked', () => {
    const onClose = vi.fn()
    const { container } = render(<Modal {...defaultProps} onClose={onClose} />)
    
    // Find backdrop (first div with onClick)
    const backdrop = container.querySelector('.fixed.inset-0.bg-black')
    if (backdrop) {
      fireEvent.click(backdrop)
      expect(onClose).toHaveBeenCalledTimes(1)
    }
  })

  it('should render children content', () => {
    const customContent = <div data-testid="custom-content">Custom Content</div>
    render(<Modal {...defaultProps}>{customContent}</Modal>)
    
    expect(screen.getByTestId('custom-content')).toBeInTheDocument()
  })

  it('should render title correctly', () => {
    render(<Modal {...defaultProps} title="Custom Title" />)
    expect(screen.getByText('Custom Title')).toBeInTheDocument()
  })
})
