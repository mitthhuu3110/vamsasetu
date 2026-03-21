import React from 'react';
import { motion } from 'framer-motion';

export interface CardProps {
  children: React.ReactNode;
  className?: string;
  variant?: 'default' | 'elevated' | 'outlined';
  padding?: 'none' | 'sm' | 'md' | 'lg';
  hoverable?: boolean;
  onClick?: () => void;
}

const Card: React.FC<CardProps> = ({
  children,
  className = '',
  variant = 'default',
  padding = 'md',
  hoverable = false,
  onClick,
}) => {
  const baseStyles = 'bg-white rounded-lg transition-all duration-200';
  
  const variantStyles = {
    default: 'shadow-sm border border-gray-200',
    elevated: 'shadow-lg',
    outlined: 'border-2 border-saffron',
  };
  
  const paddingStyles = {
    none: '',
    sm: 'p-3',
    md: 'p-4',
    lg: 'p-6',
  };
  
  const interactiveStyles = hoverable || onClick
    ? 'cursor-pointer hover:shadow-xl'
    : '';
  
  const cardClasses = `${baseStyles} ${variantStyles[variant]} ${paddingStyles[padding]} ${interactiveStyles} ${className}`;
  
  const cardProps = {
    className: cardClasses,
    onClick,
    ...(onClick && {
      role: 'button',
      tabIndex: 0,
      onKeyDown: (e: React.KeyboardEvent) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          onClick();
        }
      },
    }),
  };
  
  if (hoverable || onClick) {
    return (
      <motion.div
        {...cardProps}
        whileHover={{ scale: 1.02, y: -2 }}
        whileTap={{ scale: 0.98 }}
      >
        {children}
      </motion.div>
    );
  }
  
  return <div {...cardProps}>{children}</div>;
};

// Card sub-components for better composition
export const CardHeader: React.FC<{
  children: React.ReactNode;
  className?: string;
}> = ({ children, className = '' }) => (
  <div className={`mb-4 ${className}`}>{children}</div>
);

export const CardTitle: React.FC<{
  children: React.ReactNode;
  className?: string;
}> = ({ children, className = '' }) => (
  <h3 className={`text-xl font-semibold text-charcoal font-heading ${className}`}>
    {children}
  </h3>
);

export const CardDescription: React.FC<{
  children: React.ReactNode;
  className?: string;
}> = ({ children, className = '' }) => (
  <p className={`text-sm text-gray-600 mt-1 ${className}`}>{children}</p>
);

export const CardContent: React.FC<{
  children: React.ReactNode;
  className?: string;
}> = ({ children, className = '' }) => (
  <div className={className}>{children}</div>
);

export const CardFooter: React.FC<{
  children: React.ReactNode;
  className?: string;
}> = ({ children, className = '' }) => (
  <div className={`mt-4 pt-4 border-t border-gray-200 ${className}`}>
    {children}
  </div>
);

export default Card;
