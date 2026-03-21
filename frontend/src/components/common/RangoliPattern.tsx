import React from 'react';

/**
 * RangoliPattern Component
 * 
 * A decorative SVG background component inspired by traditional Indian rangoli designs.
 * Features geometric patterns with symmetry using VamsaSetu theme colors.
 * 
 * **Validates: Requirements 10.3**
 * 
 * @example
 * ```tsx
 * // As a background overlay
 * <div className="relative">
 *   <RangoliPattern opacity={0.1} />
 *   <div className="relative z-10">Content here</div>
 * </div>
 * 
 * // With custom size
 * <RangoliPattern size={400} opacity={0.15} />
 * ```
 */

interface RangoliPatternProps {
  /** Opacity of the pattern (0-1), default 0.1 for subtle background */
  opacity?: number;
  /** Size of the pattern tile in pixels, default 200 */
  size?: number;
  /** Additional CSS classes */
  className?: string;
}

export const RangoliPattern: React.FC<RangoliPatternProps> = ({
  opacity = 0.1,
  size = 200,
  className = '',
}) => {
  const patternId = `rangoli-pattern-${Math.random().toString(36).substr(2, 9)}`;
  
  return (
    <div 
      className={`absolute inset-0 pointer-events-none ${className}`}
      style={{ opacity }}
      aria-hidden="true"
    >
      <svg
        width="100%"
        height="100%"
        xmlns="http://www.w3.org/2000/svg"
      >
        <defs>
          <pattern
            id={patternId}
            x="0"
            y="0"
            width={size}
            height={size}
            patternUnits="userSpaceOnUse"
          >
            {/* Central lotus-inspired mandala */}
            <g transform={`translate(${size / 2}, ${size / 2})`}>
              {/* Outer petals - Saffron */}
              {[0, 45, 90, 135, 180, 225, 270, 315].map((angle) => (
                <ellipse
                  key={`petal-outer-${angle}`}
                  cx="0"
                  cy="-40"
                  rx="15"
                  ry="35"
                  fill="#E8650A"
                  transform={`rotate(${angle})`}
                />
              ))}
              
              {/* Middle petals - Turmeric */}
              {[22.5, 67.5, 112.5, 157.5, 202.5, 247.5, 292.5, 337.5].map((angle) => (
                <ellipse
                  key={`petal-middle-${angle}`}
                  cx="0"
                  cy="-28"
                  rx="12"
                  ry="25"
                  fill="#F5A623"
                  transform={`rotate(${angle})`}
                />
              ))}
              
              {/* Inner circle - Teal */}
              <circle cx="0" cy="0" r="20" fill="#0D4A52" />
              
              {/* Center dot - Saffron */}
              <circle cx="0" cy="0" r="8" fill="#E8650A" />
              
              {/* Decorative dots around center - Turmeric */}
              {[0, 60, 120, 180, 240, 300].map((angle) => {
                const rad = (angle * Math.PI) / 180;
                const x = Math.cos(rad) * 12;
                const y = Math.sin(rad) * 12;
                return (
                  <circle
                    key={`dot-${angle}`}
                    cx={x}
                    cy={y}
                    r="2"
                    fill="#F5A623"
                  />
                );
              })}
            </g>
            
            {/* Corner decorative elements - Teal */}
            {[
              [0, 0],
              [size, 0],
              [0, size],
              [size, size],
            ].map(([x, y], idx) => (
              <g key={`corner-${idx}`} transform={`translate(${x}, ${y})`}>
                <circle cx="0" cy="0" r="8" fill="#0D4A52" />
                <circle cx="0" cy="0" r="3" fill="#F5A623" />
              </g>
            ))}
            
            {/* Connecting lines - subtle geometric grid */}
            <line
              x1={size / 2}
              y1="0"
              x2={size / 2}
              y2={size}
              stroke="#0D4A52"
              strokeWidth="0.5"
              opacity="0.3"
            />
            <line
              x1="0"
              y1={size / 2}
              x2={size}
              y2={size / 2}
              stroke="#0D4A52"
              strokeWidth="0.5"
              opacity="0.3"
            />
          </pattern>
        </defs>
        
        <rect width="100%" height="100%" fill={`url(#${patternId})`} />
      </svg>
    </div>
  );
};

export default RangoliPattern;
