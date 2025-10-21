import { motion } from 'motion/react';
import { LucideIcon } from 'lucide-react';
import { Card } from './ui/card';

interface StatsCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  gradient: string;
  delay?: number;
}

export function StatsCard({ title, value, icon: Icon, gradient, delay = 0 }: StatsCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, delay }}
    >
      <Card className="relative overflow-hidden p-6 border-border bg-card shadow-lg hover:shadow-xl transition-all duration-300">
        {/* Background Gradient Orb */}
        <div
          className={`absolute -right-4 -top-4 w-24 h-24 rounded-full opacity-10 blur-2xl ${gradient}`}
        />
        
        <div className="relative flex items-start justify-between">
          <div className="space-y-2">
            <p className="text-muted-foreground text-sm">{title}</p>
            <p className="text-3xl tracking-tight">{value}</p>
          </div>
          
          <div className={`p-3 rounded-xl ${gradient} bg-opacity-20`}>
            <Icon className="w-6 h-6 text-primary" />
          </div>
        </div>
      </Card>
    </motion.div>
  );
}
