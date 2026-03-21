import React from 'react';
import Navbar from './Navbar';
import Sidebar from './Sidebar';
import BottomNav from './BottomNav';
import { useResponsive } from '../../hooks/useResponsive';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { isMobile } = useResponsive();

  return (
    <div className="min-h-screen bg-ivory">
      <Navbar />
      <div className="flex">
        {/* Show Sidebar only on desktop */}
        {!isMobile && <Sidebar />}
        <main className="flex-1 p-4 md:p-6 pb-20 md:pb-6">
          {children}
        </main>
      </div>
      {/* Show BottomNav only on mobile */}
      {isMobile && <BottomNav />}
    </div>
  );
};

export default Layout;
