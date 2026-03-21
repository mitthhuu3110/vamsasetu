/**
 * UI Components Usage Examples
 * 
 * This file demonstrates how to use the VamsaSetu UI components
 * with proper TypeScript types and accessibility features.
 */

import React, { useState } from 'react';
import Button from './Button';
import Input from './Input';
import Modal from './Modal';
import Card, { CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from './Card';

export const ButtonExamples: React.FC = () => {
  return (
    <div className="space-y-4 p-6">
      <h2 className="text-2xl font-heading font-semibold">Button Examples</h2>
      
      {/* Primary Button */}
      <Button variant="primary" onClick={() => alert('Primary clicked')}>
        Primary Button
      </Button>
      
      {/* Secondary Button */}
      <Button variant="secondary" onClick={() => alert('Secondary clicked')}>
        Secondary Button
      </Button>
      
      {/* Outline Button */}
      <Button variant="outline" onClick={() => alert('Outline clicked')}>
        Outline Button
      </Button>
      
      {/* Button with Icons */}
      <Button
        variant="primary"
        leftIcon={<span>👤</span>}
        onClick={() => alert('With icon')}
      >
        Add Member
      </Button>
      
      {/* Loading Button */}
      <Button variant="primary" isLoading>
        Saving...
      </Button>
      
      {/* Disabled Button */}
      <Button variant="primary" disabled>
        Disabled Button
      </Button>
      
      {/* Full Width Button */}
      <Button variant="primary" fullWidth>
        Full Width Button
      </Button>
      
      {/* Different Sizes */}
      <div className="flex gap-2 items-center">
        <Button size="sm">Small</Button>
        <Button size="md">Medium</Button>
        <Button size="lg">Large</Button>
      </div>
    </div>
  );
};

export const InputExamples: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [emailError, setEmailError] = useState('');
  
  const validateEmail = (value: string) => {
    if (!value) {
      setEmailError('Email is required');
    } else if (!/\S+@\S+\.\S+/.test(value)) {
      setEmailError('Invalid email format');
    } else {
      setEmailError('');
    }
  };
  
  return (
    <div className="space-y-4 p-6 max-w-md">
      <h2 className="text-2xl font-heading font-semibold">Input Examples</h2>
      
      {/* Basic Input */}
      <Input
        label="Username"
        placeholder="Enter your username"
      />
      
      {/* Required Input */}
      <Input
        label="Email"
        type="email"
        placeholder="your@email.com"
        required
        value={email}
        onChange={(e) => {
          setEmail(e.target.value);
          validateEmail(e.target.value);
        }}
        error={emailError}
      />
      
      {/* Input with Helper Text */}
      <Input
        label="Password"
        type="password"
        placeholder="Enter password"
        helperText="Must be at least 8 characters"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      
      {/* Input with Left Icon */}
      <Input
        label="Search"
        placeholder="Search members..."
        leftIcon={
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        }
      />
      
      {/* Full Width Input */}
      <Input
        label="Full Name"
        placeholder="Enter full name"
        fullWidth
      />
    </div>
  );
};

export const ModalExamples: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  
  return (
    <div className="space-y-4 p-6">
      <h2 className="text-2xl font-heading font-semibold">Modal Examples</h2>
      
      {/* Basic Modal */}
      <Button onClick={() => setIsOpen(true)}>
        Open Modal
      </Button>
      
      <Modal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Add New Member"
        size="md"
      >
        <div className="space-y-4">
          <Input label="Name" placeholder="Enter member name" fullWidth />
          <Input label="Date of Birth" type="date" fullWidth />
          <div className="flex gap-2 justify-end mt-6">
            <Button variant="outline" onClick={() => setIsOpen(false)}>
              Cancel
            </Button>
            <Button variant="primary" onClick={() => setIsOpen(false)}>
              Save Member
            </Button>
          </div>
        </div>
      </Modal>
      
      {/* Confirmation Modal */}
      <Button variant="secondary" onClick={() => setConfirmOpen(true)}>
        Delete Confirmation
      </Button>
      
      <Modal
        isOpen={confirmOpen}
        onClose={() => setConfirmOpen(false)}
        title="Confirm Deletion"
        size="sm"
        closeOnOverlayClick={false}
      >
        <div className="space-y-4">
          <p className="text-gray-600">
            Are you sure you want to delete this member? This action cannot be undone.
          </p>
          <div className="flex gap-2 justify-end">
            <Button variant="outline" onClick={() => setConfirmOpen(false)}>
              Cancel
            </Button>
            <Button
              variant="primary"
              className="bg-rose hover:bg-rose"
              onClick={() => setConfirmOpen(false)}
            >
              Delete
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};

export const CardExamples: React.FC = () => {
  return (
    <div className="space-y-4 p-6">
      <h2 className="text-2xl font-heading font-semibold">Card Examples</h2>
      
      {/* Basic Card */}
      <Card>
        <p>This is a basic card with default styling.</p>
      </Card>
      
      {/* Card with Header and Footer */}
      <Card variant="elevated">
        <CardHeader>
          <CardTitle>Family Member</CardTitle>
          <CardDescription>Details about the family member</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-gray-600">
            Name: Rajesh Kumar<br />
            Date of Birth: January 15, 1980<br />
            Relationship: Father
          </p>
        </CardContent>
        <CardFooter>
          <div className="flex gap-2">
            <Button size="sm" variant="outline">Edit</Button>
            <Button size="sm" variant="primary">View Details</Button>
          </div>
        </CardFooter>
      </Card>
      
      {/* Outlined Card */}
      <Card variant="outlined" padding="lg">
        <CardTitle>Important Notice</CardTitle>
        <CardContent className="mt-2">
          <p className="text-gray-600">
            This card uses the outlined variant with saffron border.
          </p>
        </CardContent>
      </Card>
      
      {/* Hoverable Card */}
      <Card hoverable onClick={() => alert('Card clicked')}>
        <CardTitle>Clickable Card</CardTitle>
        <CardContent className="mt-2">
          <p className="text-gray-600">
            Hover over this card to see the animation effect.
          </p>
        </CardContent>
      </Card>
      
      {/* Grid of Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card hoverable>
          <CardTitle>Card 1</CardTitle>
          <CardContent className="mt-2">
            <p className="text-sm text-gray-600">Content for card 1</p>
          </CardContent>
        </Card>
        <Card hoverable>
          <CardTitle>Card 2</CardTitle>
          <CardContent className="mt-2">
            <p className="text-sm text-gray-600">Content for card 2</p>
          </CardContent>
        </Card>
        <Card hoverable>
          <CardTitle>Card 3</CardTitle>
          <CardContent className="mt-2">
            <p className="text-sm text-gray-600">Content for card 3</p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

// Complete Example: Form in Modal with Card
export const CompleteExample: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
  });
  const [errors, setErrors] = useState({
    name: '',
    email: '',
    phone: '',
  });
  
  const handleSubmit = () => {
    // Validation
    const newErrors = {
      name: formData.name ? '' : 'Name is required',
      email: formData.email ? '' : 'Email is required',
      phone: formData.phone ? '' : 'Phone is required',
    };
    
    setErrors(newErrors);
    
    if (!newErrors.name && !newErrors.email && !newErrors.phone) {
      alert('Form submitted successfully!');
      setIsModalOpen(false);
      setFormData({ name: '', email: '', phone: '' });
    }
  };
  
  return (
    <div className="p-6">
      <Card variant="elevated" padding="lg">
        <CardHeader>
          <CardTitle>Member Management</CardTitle>
          <CardDescription>Add and manage family members</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-gray-600 mb-4">
            Click the button below to add a new family member to your tree.
          </p>
          <Button variant="primary" onClick={() => setIsModalOpen(true)}>
            Add New Member
          </Button>
        </CardContent>
      </Card>
      
      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Add Family Member"
        size="md"
      >
        <div className="space-y-4">
          <Input
            label="Full Name"
            placeholder="Enter full name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            error={errors.name}
            required
            fullWidth
          />
          
          <Input
            label="Email"
            type="email"
            placeholder="email@example.com"
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            error={errors.email}
            required
            fullWidth
          />
          
          <Input
            label="Phone"
            type="tel"
            placeholder="+91 98765 43210"
            value={formData.phone}
            onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
            error={errors.phone}
            helperText="Include country code"
            required
            fullWidth
          />
          
          <div className="flex gap-2 justify-end pt-4">
            <Button
              variant="outline"
              onClick={() => {
                setIsModalOpen(false);
                setFormData({ name: '', email: '', phone: '' });
                setErrors({ name: '', email: '', phone: '' });
              }}
            >
              Cancel
            </Button>
            <Button variant="primary" onClick={handleSubmit}>
              Add Member
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};
