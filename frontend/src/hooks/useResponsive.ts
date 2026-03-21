import { useState, useEffect } from 'react';

export type Breakpoint = 'mobile' | 'tablet' | 'desktop';

export interface ResponsiveState {
  breakpoint: Breakpoint;
  isMobile: boolean;
  isTablet: boolean;
  isDesktop: boolean;
  width: number;
}

/**
 * useResponsive Hook
 * 
 * Detects current responsive breakpoint and provides boolean helpers.
 * 
 * Breakpoints:
 * - mobile: < 768px
 * - tablet: 768px - 1024px
 * - desktop: > 1024px
 * 
 * @returns ResponsiveState with breakpoint info and boolean helpers
 * 
 * @example
 * ```tsx
 * const { isMobile, isDesktop, breakpoint } = useResponsive();
 * 
 * return (
 *   <div>
 *     {isMobile && <MobileView />}
 *     {isDesktop && <DesktopView />}
 *     <p>Current breakpoint: {breakpoint}</p>
 *   </div>
 * );
 * ```
 * 
 * **Validates: Requirements 9.4**
 */
export const useResponsive = (): ResponsiveState => {
  const [state, setState] = useState<ResponsiveState>(() => {
    const width = typeof window !== 'undefined' ? window.innerWidth : 1024;
    return getResponsiveState(width);
  });

  useEffect(() => {
    const handleResize = () => {
      setState(getResponsiveState(window.innerWidth));
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return state;
};

function getResponsiveState(width: number): ResponsiveState {
  let breakpoint: Breakpoint;
  
  if (width < 768) {
    breakpoint = 'mobile';
  } else if (width < 1024) {
    breakpoint = 'tablet';
  } else {
    breakpoint = 'desktop';
  }

  return {
    breakpoint,
    isMobile: breakpoint === 'mobile',
    isTablet: breakpoint === 'tablet',
    isDesktop: breakpoint === 'desktop',
    width,
  };
}

export default useResponsive;
