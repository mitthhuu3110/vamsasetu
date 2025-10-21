import { motion } from 'motion/react';
import { Calendar, Heart, Sparkles } from 'lucide-react';
import { Card } from './ui/card';
import { Button } from './ui/button';
import { Badge } from './ui/badge';

interface EventCardProps {
  name: string;
  relation: string;
  date: string;
  type: 'birthday' | 'anniversary' | 'puja';
  daysLeft: number;
  avatarUrl?: string;
}

export function EventCard({ name, relation, date, type, daysLeft, avatarUrl }: EventCardProps) {
  const getEventIcon = () => {
    switch (type) {
      case 'birthday':
        return <Heart className="w-4 h-4" />;
      case 'anniversary':
        return <Sparkles className="w-4 h-4" />;
      case 'puja':
        return <Calendar className="w-4 h-4" />;
    }
  };

  const getEventColor = () => {
    switch (type) {
      case 'birthday':
        return 'bg-rose-100 text-rose-700 border-rose-200';
      case 'anniversary':
        return 'bg-amber-100 text-amber-700 border-amber-200';
      case 'puja':
        return 'bg-emerald-100 text-emerald-700 border-emerald-200';
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      whileHover={{ scale: 1.02, y: -4 }}
      transition={{ duration: 0.3 }}
    >
      <Card className="relative overflow-hidden bg-card border-border shadow-md hover:shadow-xl transition-all duration-300">
        {/* Decorative Pattern */}
        <div className="absolute top-0 right-0 w-32 h-32 opacity-5">
          <svg viewBox="0 0 100 100" className="w-full h-full">
            <circle cx="50" cy="50" r="40" fill="none" stroke="currentColor" strokeWidth="2" />
            <circle cx="50" cy="50" r="30" fill="none" stroke="currentColor" strokeWidth="2" />
            <circle cx="50" cy="50" r="20" fill="none" stroke="currentColor" strokeWidth="2" />
          </svg>
        </div>

        <div className="p-5 relative">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-start space-x-3">
              {/* Avatar */}
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center overflow-hidden shadow-md">
                {avatarUrl ? (
                  <img src={avatarUrl} alt={name} className="w-full h-full object-cover" />
                ) : (
                  <span className="text-primary-foreground">
                    {name.charAt(0)}
                  </span>
                )}
              </div>

              {/* Event Info */}
              <div className="flex-1">
                <h3 className="text-foreground mb-1">{name}</h3>
                <p className="text-sm text-muted-foreground">{relation}</p>
              </div>
            </div>

            {/* Event Type Badge */}
            <Badge className={`${getEventColor()} border flex items-center space-x-1`}>
              {getEventIcon()}
              <span className="capitalize">{type}</span>
            </Badge>
          </div>

          {/* Date Info */}
          <div className="flex items-center justify-between pt-3 border-t border-border">
            <div className="flex items-center space-x-2 text-sm">
              <Calendar className="w-4 h-4 text-muted-foreground" />
              <span className="text-muted-foreground">{date}</span>
            </div>
            
            {daysLeft <= 7 && (
              <span className="text-xs px-2 py-1 rounded-full bg-primary/20 text-primary">
                {daysLeft === 0 ? 'Today!' : `${daysLeft} days left`}
              </span>
            )}
          </div>

          {/* Send Wishes Button */}
          <Button className="w-full mt-4 bg-gradient-to-r from-primary to-secondary hover:opacity-90 text-primary-foreground shadow-md">
            Send Wishes
          </Button>
        </div>
      </Card>
    </motion.div>
  );
}
