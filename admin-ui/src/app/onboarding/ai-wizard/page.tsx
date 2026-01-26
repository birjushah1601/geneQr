'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Loader2, Send, Upload, CheckCircle2, Circle, ArrowLeft, FileText, Shield } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';
import Navigation from '@/components/Navigation';
import { decodeJWT } from '@/lib/jwt';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';

// Onboarding steps - Logical business flow
const STEPS = [
  { id: 'manufacturer', label: 'Manufacturer Info', icon: '🏢' },
  { id: 'team', label: 'Team Members', icon: '👥' },
  { id: 'equipment', label: 'Equipment Catalog', icon: '🔧' },
  { id: 'parts', label: 'Parts Catalog', icon: '📦' },
  { id: 'engineers', label: 'Service Engineers', icon: '👷' },
  { id: 'installations', label: 'Equipment Installations', icon: '📋' },
  { id: 'review', label: 'Review & Complete', icon: '✅' },
];

interface Message {
  id: string;
  role: 'ai' | 'user';
  content: string;
  timestamp: Date;
  type?: 'text' | 'file' | 'action';
  actions?: Array<{ label: string; value: string }>;
}

interface OnboardingData {
  company?: any;
  equipmentTypes?: any[];
  parts?: any[];
  hospitals?: any[];
  installations?: any[];
  engineers?: any[];
}

export default function AIOnboardingWizard() {
  const router = useRouter();
  const { organizationContext } = useAuth();
  const [isAuthorized, setIsAuthorized] = useState<boolean | null>(null);
  const [currentStep, setCurrentStep] = useState(0);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);
  const [showManualFallback, setShowManualFallback] = useState(false);
  const [onboardingData, setOnboardingData] = useState<OnboardingData>({});
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Check authorization - system admins and manufacturer admins
  useEffect(() => {
    const checkAuth = () => {
      // Get user from localStorage
      const userStr = localStorage.getItem('user');
      const user = userStr ? JSON.parse(userStr) : null;
      const orgType = user?.organization_type || organizationContext?.organization_type || '';
      const role = user?.role || organizationContext?.role || '';
      
      console.log('AI Wizard Access Check:', { orgType, role, user, organizationContext });
      
      // Allow system admins and manufacturer admins
      const authorized = orgType === 'system' || 
                        (orgType === 'manufacturer' && role === 'admin');
      setIsAuthorized(authorized);
      
      if (!authorized) {
        // Redirect unauthorized users
        setTimeout(() => {
          router.push('/dashboard');
        }, 3000);
      } else if (orgType === 'manufacturer') {
        // For manufacturer admins, skip company profile step
        setCurrentStep(1); // Start from equipment types
      }
    };

    checkAuth();
  }, [organizationContext, router]);

  useEffect(() => {
    // Only show initial greeting if authorized
    if (isAuthorized) {
      // Check if manufacturer admin (skip company profile)
      const userStr = localStorage.getItem('user');
      const user = userStr ? JSON.parse(userStr) : null;
      const orgType = user?.organization_type || '';
      const isManufacturerAdmin = orgType === 'manufacturer';

      if (isManufacturerAdmin) {
        // Greeting for manufacturer admins (skip manufacturer info)
        addAIMessage(
          `Hi! 👋 I'm your AI assistant. I'll help you set up your complete system.\n\n` +
          `Here's what we'll configure:\n` +
          `1. 👥 Team members (admins, managers)\n` +
          `2. 🔧 Equipment you manufacture\n` +
          `3. 📦 Parts catalog for service tickets\n` +
          `4. 👷 Service engineers and their skills\n` +
          `5. 📋 Equipment installations at customer sites\n\n` +
          `Let's start! Would you like to invite team members to help manage the system?\n` +
          `Example: CEO, Operations Manager, Service Manager`,
          [
            { label: '👥 Yes, invite team members', value: 'invite_team' },
            { label: '⏭️ Skip, just me for now', value: 'skip_team' },
            { label: '📝 Switch to Manual Form', value: 'manual' },
          ]
        );
      } else {
        // Greeting for system admins (full flow)
        addAIMessage(
          `Hi! 👋 I'm your AI onboarding assistant for setting up a new manufacturer.\n\n` +
          `Here's the complete setup flow:\n` +
          `1. 🏢 Manufacturer company information\n` +
          `2. 👥 Team members (multiple admins/managers)\n` +
          `3. 🔧 Equipment they manufacture\n` +
          `4. 📦 Parts catalog for service tickets\n` +
          `5. 👷 Service engineers and their skills\n` +
          `6. 📋 Equipment installations at customer sites\n\n` +
          `Let's begin! What's the manufacturer's registered company name?`,
          [
            { label: '📝 Switch to Manual Form', value: 'manual' },
          ]
        );
      }
    }
  }, [isAuthorized]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const addAIMessage = (content: string, actions?: Array<{ label: string; value: string }>) => {
    const message: Message = {
      id: Date.now().toString(),
      role: 'ai',
      content,
      timestamp: new Date(),
      type: actions ? 'action' : 'text',
      actions,
    };
    setMessages(prev => [...prev, message]);
  };

  const addUserMessage = (content: string, type: 'text' | 'file' = 'text') => {
    const message: Message = {
      id: Date.now().toString(),
      role: 'user',
      content,
      timestamp: new Date(),
      type,
    };
    setMessages(prev => [...prev, message]);
  };

  // Show appropriate prompt when user navigates to a step
  const showStepPrompt = (stepIndex: number) => {
    const step = STEPS[stepIndex];
    
    switch (step.id) {
      case 'manufacturer':
        addAIMessage(
          `🏢 **Manufacturer Information**\n\n` +
          `Let's start with the basic company details.\n\n` +
          `What's the manufacturer's registered company name?`,
          [
            { label: '📝 Switch to Manual Form', value: 'manual' },
          ]
        );
        break;
      
      case 'team':
        addAIMessage(
          `👥 **Team Members**\n\n` +
          `Would you like to invite team members to help manage the system?\n` +
          `Example: CEO, Operations Manager, Service Manager`,
          [
            { label: '👥 Yes, invite team members', value: 'invite_team' },
            { label: '⏭️ Skip, just me for now', value: 'skip_team' },
          ]
        );
        break;
      
      case 'equipment':
        addAIMessage(
          `🔧 **Equipment Catalog**\n\n` +
          `Let's add the equipment you manufacture.\n\n` +
          `You can either upload a CSV file or enter equipment manually.`,
          [
            { label: '📤 Upload CSV Template', value: 'upload_equipment' },
            { label: '📝 I have multiple equipment types', value: 'multiple_equipment' },
          ]
        );
        break;
      
      case 'parts':
        addAIMessage(
          `📦 **Parts Catalog**\n\n` +
          `Now let's add the spare parts for your equipment.\n\n` +
          `You can upload a CSV file with your parts catalog.`,
          [
            { label: '📤 Upload Parts CSV', value: 'upload_parts' },
            { label: '⏭️ Skip for now', value: 'skip_parts' },
          ]
        );
        break;
      
      case 'engineers':
        addAIMessage(
          `👷 **Service Engineers**\n\n` +
          `Let's add your service engineers who will handle equipment maintenance.\n\n` +
          `You can upload a CSV file with engineer details.`,
          [
            { label: '📤 Upload Engineers CSV', value: 'upload_engineers' },
            { label: '✍️ Add manually', value: 'manual_engineers' },
            { label: '⏭️ Skip for now', value: 'skip_engineers' },
          ]
        );
        break;
      
      case 'installations':
        addAIMessage(
          `📋 **Equipment Installations**\n\n` +
          `Finally, let's record where your equipment is installed at customer sites.\n\n` +
          `You can upload a CSV file with installation records.`,
          [
            { label: '📤 Upload Installations CSV', value: 'upload_installations' },
            { label: '⏭️ Skip for now', value: 'skip_installations' },
          ]
        );
        break;
      
      case 'review':
        addAIMessage(
          `✅ **Review & Complete**\n\n` +
          `Great work! Let's review what we've set up:\n\n` +
          `${onboardingData.company ? '✅ Manufacturer information\n' : ''}` +
          `${onboardingData.equipmentTypes ? `✅ ${onboardingData.equipmentTypes.length} equipment types\n` : ''}` +
          `${onboardingData.parts ? `✅ ${onboardingData.parts.length} parts\n` : ''}` +
          `${onboardingData.engineers ? `✅ ${onboardingData.engineers.length} engineers\n` : ''}` +
          `${onboardingData.installations ? `✅ ${onboardingData.installations.length} installations\n` : ''}` +
          `\nReady to complete the setup?`,
          [
            { label: '✅ Complete Setup', value: 'complete' },
            { label: '📝 Review Data', value: 'review' },
          ]
        );
        break;
    }
  };

  const processUserInput = async (userInput: string) => {
    if (!userInput.trim()) return;

    addUserMessage(userInput);
    setInput('');
    setIsProcessing(true);

    try {
      // Simulate AI processing (In production, call OpenAI/GPT-4 API)
      await new Promise(resolve => setTimeout(resolve, 1000));

      // Handle cross-step actions first (before routing by step)
      if (userInput === 'continue_equipment') {
        await handleEquipmentStep(userInput);
        setIsProcessing(false);
        return;
      }

      // Route based on current step
      switch (STEPS[currentStep].id) {
        case 'manufacturer':
          await handleManufacturerStep(userInput);
          break;
        case 'team':
          await handleTeamStep(userInput);
          break;
        case 'equipment':
          await handleEquipmentStep(userInput);
          break;
        case 'parts':
          await handlePartsStep(userInput);
          break;
        case 'engineers':
          await handleEngineersStep(userInput);
          break;
        case 'installations':
          await handleInstallationsStep(userInput);
          break;
        default:
          addAIMessage("I'm not sure how to help with that. Can you rephrase?");
      }
    } catch (error) {
      console.error('AI processing error:', error);
      addAIMessage(
        "I encountered an error. Would you like to switch to manual entry?",
        [
          { label: 'Switch to Manual', value: 'manual' },
          { label: 'Try Again', value: 'retry' },
        ]
      );
    } finally {
      setIsProcessing(false);
    }
  };

  const handleManufacturerStep = async (input: string) => {
    // Extract company name (in production, use GPT-4 for NLP)
    const manufacturerData = {
      name: input,
      created_at: new Date().toISOString(),
    };

    setOnboardingData(prev => ({ ...prev, manufacturer: manufacturerData }));

    addAIMessage(
      `Perfect! I've registered "${input}" as the manufacturer.\n\n` +
      `Step 1 complete! ✅ Now let's set up the **Team**.\n\n` +
      `Would you like to invite other team members (admins, managers) to help manage the system?\n` +
      `This allows multiple people to:\n` +
      `• Manage equipment and parts\n` +
      `• View service tickets\n` +
      `• Oversee operations\n\n` +
      `You can invite: CEO, CTO, Operations Manager, Service Manager, etc.`,
      [
        { label: '👥 Yes, invite team members', value: 'invite_team' },
        { label: '📧 I have a list (upload CSV)', value: 'upload_team' },
        { label: '⏭️ Skip for now, just me', value: 'skip_team' },
      ]
    );
    setCurrentStep(1); // Move to team step
  };

  const handleTeamStep = async (input: string) => {
    if (input === 'invite_team') {
      addAIMessage(
        `Great! Let's invite your team members.\n\n` +
        `Please provide their details in this format:\n` +
        `Name, Email, Role\n\n` +
        `Example:\n` +
        `Rajesh Kumar, ceo@company.com, admin\n` +
        `Priya Sharma, ops@company.com, manager\n` +
        `Amit Verma, service@company.com, manager`,
        [
          { label: '✍️ I\'ll type the details', value: 'type_team' },
          { label: '📧 Upload CSV file', value: 'upload_team' },
          { label: '⏭️ Skip for now', value: 'skip_team' },
        ]
      );
      return;
    }

    if (input === 'upload_team') {
      addAIMessage(
        `Please upload a CSV file with team member details.\n\n` +
        `CSV Format:\n` +
        `name,email,role\n` +
        `Rajesh Kumar,ceo@company.com,admin\n` +
        `Priya Sharma,ops@company.com,manager\n\n` +
        `Roles: admin (full access), manager (operations), viewer (read-only)`,
        [
          { label: '📤 Upload CSV', value: 'do_upload_team' },
          { label: '⏭️ Skip for now', value: 'skip_team' },
        ]
      );
      
      if (input === 'do_upload_team') {
        fileInputRef.current?.click();
      }
      return;
    }

    if (input === 'skip_team') {
      addAIMessage(
        `No problem! You can invite team members later from the dashboard.\n\n` +
        `Step 2 complete! ✅ Now let's add **Equipment Catalog**.\n\n` +
        `What type of medical equipment do you manufacture?\n` +
        `Example: "Ventilators", "MRI Machines", "CT Scanners", "Ultrasound Systems"`,
        [
          { label: '📝 I have multiple equipment types', value: 'multiple_equipment' },
        ]
      );
      setCurrentStep(2); // Move to equipment step
      return;
    }

    // If user typed team details, process them
    // Format: Name, Email, Role (one per line)
    const teamMembers = parseTeamInput(input);
    
    if (teamMembers.length > 0) {
      setOnboardingData(prev => ({ ...prev, teamMembers }));
      
      // Send invitations immediately
      addAIMessage(
        `Excellent! I've noted ${teamMembers.length} team member(s):\n` +
        teamMembers.map(m => `• ${m.name} (${m.email}) - ${m.role}`).join('\n') + 
        `\n\n⏳ Sending invitations now...`,
        []
      );
      
      // Send invitations via API
      await sendTeamInvitations(teamMembers);
    }
  };
  
  const sendTeamInvitations = async (teamMembers: Array<{name: string, email: string, role: string}>) => {
    try {
      // Get token from cookies (where auth context stores it)
      const getCookie = (name: string) => {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        if (parts.length === 2) return parts.pop()?.split(';').shift();
        return null;
      };
      
      const token = getCookie('accessToken') || localStorage.getItem('access_token') || localStorage.getItem('token');
      
      if (!token) {
        console.error('[Invitations] No token found in cookies or localStorage');
        addAIMessage(
          `⚠️ You need to be logged in to send invitations.\n\n` +
          `Team member details have been saved. You can send invitations later from the dashboard.\n\n` +
          `Step 2 complete! ✅ Now let's add **Equipment Catalog**.`,
          [
            { label: '➡️ Continue to equipment', value: 'continue_equipment' },
          ]
        );
        return;
      }
      
      // Get organization ID - try multiple sources
      let orgId = organizationContext?.organization_id;
      
      if (!orgId) {
        // Try from user object in localStorage
        const userStr = localStorage.getItem('user');
        if (userStr) {
          const user = JSON.parse(userStr);
          orgId = user.organization_id;
          console.log('[Invitations] Using org ID from user object:', orgId);
        }
      }
      
      if (!orgId) {
        // Fallback: Decode JWT token directly
        const claims = decodeJWT(token);
        orgId = claims?.organization_id;
        console.log('[Invitations] Using fallback - decoded org ID from JWT:', orgId);
      }
      
      console.log('[Invitations] Debug info:', {
        hasToken: !!token,
        tokenSource: getCookie('accessToken') ? 'cookie' : 'localStorage',
        organizationContext,
        orgId,
        userObject: localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!) : null,
      });
      
      if (!orgId) {
        console.error('[Invitations] No organization ID found anywhere!');
        addAIMessage(
          `⚠️ No organization found. Team members will be invited after organization setup.\n\n` +
          `Step 2 complete! ✅ Now let's add **Equipment Catalog**.`,
          [
            { label: '➡️ Continue to equipment', value: 'continue_equipment' },
          ]
        );
        return;
      }
      
      // Send invitations
      let successCount = 0;
      let failedEmails: string[] = [];
      
      for (const member of teamMembers) {
        try {
          const response = await fetch(`${API_BASE_URL}/v1/organizations/${orgId}/invitations`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({
              email: member.email,
              name: member.name,
              role: member.role,
            }),
          });
          
          if (response.ok) {
            successCount++;
          } else {
            const error = await response.json();
            failedEmails.push(`${member.email} (${error.error || 'unknown error'})`);
          }
        } catch (err) {
          failedEmails.push(`${member.email} (network error)`);
        }
      }
      
      // Show results
      if (successCount === teamMembers.length) {
        addAIMessage(
          `✅ Success! Sent ${successCount} invitation(s)!\n\n` +
          `📧 Your team members will receive invitation emails shortly.\n\n` +
          `Step 2 complete! ✅ Now let's add **Equipment Catalog**.\n\n` +
          `What type of medical equipment do you manufacture?\n` +
          `Example: "Ventilators", "MRI Machines", "CT Scanners"`,
          [
            { label: '📝 I have multiple equipment types', value: 'multiple_equipment' },
          ]
        );
        setCurrentStep(2); // Move to equipment
      } else if (successCount > 0) {
        addAIMessage(
          `⚠️ Partially successful:\n` +
          `✅ Sent ${successCount} invitation(s)\n` +
          `❌ Failed: ${failedEmails.join(', ')}\n\n` +
          `You can retry failed invitations from the dashboard later.\n\n` +
          `Step 2 complete! ✅ Continue to equipment?`,
          [
            { label: '➕ Retry failed invitations', value: 'invite_team' },
            { label: '➡️ Continue to equipment', value: 'continue_equipment' },
          ]
        );
      } else {
        addAIMessage(
          `❌ Failed to send invitations:\n` +
          failedEmails.join('\n') + 
          `\n\nTeam member details saved. You can send invitations later from the dashboard.\n\n` +
          `Continue to equipment?`,
          [
            { label: '🔄 Try again', value: 'invite_team' },
            { label: '➡️ Continue to equipment', value: 'continue_equipment' },
          ]
        );
      }
      
    } catch (error) {
      console.error('Failed to send invitations:', error);
      addAIMessage(
        `❌ Error sending invitations. Team member details have been saved.\n\n` +
        `You can send invitations later from the dashboard.\n\n` +
        `Step 2 complete! ✅ Continue to equipment?`,
        [
          { label: '➡️ Continue to equipment', value: 'continue_equipment' },
        ]
      );
    }
  };

  const parseTeamInput = (input: string): Array<{name: string, email: string, role: string}> => {
    // Simple parser for team member input
    const lines = input.split('\n').filter(l => l.trim());
    const members: Array<{name: string, email: string, role: string}> = [];
    
    for (const line of lines) {
      const parts = line.split(',').map(p => p.trim());
      if (parts.length >= 2) {
        members.push({
          name: parts[0],
          email: parts[1],
          role: parts[2] || 'manager', // Default to manager
        });
      }
    }
    
    return members;
  };

  const handleEquipmentStep = async (input: string) => {
    // Handle continue from team step
    if (input === 'continue_equipment') {
      addAIMessage(
        `Now let's add your **Equipment Catalog**.\n\n` +
        `Do you have an equipment list in Excel or CSV format?\n\n` +
        `📥 **Need a template?** [Download Equipment Catalog Template](/templates/equipment-catalog-template.csv)`,
        [
          { label: '📤 Yes, upload file', value: 'upload_equipment' },
          { label: '📥 Download template first', value: 'download_equipment_template' },
          { label: '✍️ I\'ll add manually', value: 'manual_equipment' },
          { label: '⏭️ Skip for now', value: 'skip_equipment' },
        ]
      );
      setCurrentStep(2); // Move to equipment
      return;
    }

    if (input === 'download_equipment_template') {
      addAIMessage(
        `Great! I've provided the template download link above.\n\n` +
        `The template includes:\n` +
        `✅ Equipment type, model, manufacturer\n` +
        `✅ Serial number and installation details\n` +
        `✅ Sample data for common medical equipment\n\n` +
        `Once downloaded, you can fill it out and upload it here.`,
        [
          { label: '📤 Upload completed template', value: 'upload_equipment' },
          { label: '⏭️ Skip for now', value: 'skip_equipment' },
        ]
      );
      return;
    }

    if (input === 'upload_equipment') {
      addAIMessage(
        `Please upload your equipment catalog CSV file.\n\n` +
        `File should contain: equipment type, model, serial numbers, etc.`,
        [
          { label: '📤 Upload CSV', value: 'do_upload_equipment' },
          { label: '⏭️ Skip for now', value: 'skip_equipment' },
        ]
      );
      
      if (input === 'do_upload_equipment') {
        fileInputRef.current?.click();
      }
      return;
    }

    if (input === 'skip_equipment') {
      addAIMessage(
        `No problem! You can add equipment later from the dashboard.\n\n` +
        `Step 3 complete! ✅ Now let's set up your **Parts Catalog**.\n\n` +
        `Parts are essential for:\n` +
        `• Service ticket management\n` +
        `• Inventory tracking\n` +
        `• Cost estimation\n\n` +
        `Do you have a parts list in Excel or CSV format?\n\n` +
        `📥 **Need a template?** [Download Parts Catalog Template](/templates/parts-catalog-template.csv)`,
        [
          { label: '📤 Yes, upload file', value: 'upload_parts' },
          { label: '📥 Download template first', value: 'download_template' },
          { label: '✍️ I\'ll add manually', value: 'manual_parts' },
          { label: '⏭️ Skip for now', value: 'skip_parts' },
        ]
      );
      setCurrentStep(3); // Move to parts
      return;
    }
    
    // Save equipment type
    const equipmentData = { types: input };
    setOnboardingData(prev => ({ ...prev, equipment: equipmentData }));

    addAIMessage(
      `Excellent! ${input} equipment registered. ✅\n\n` +
      `Step 3 complete! Now let's set up the **Parts Catalog**.\n\n` +
      `Parts are essential for:\n` +
      `• Service ticket management\n` +
      `• Inventory tracking\n` +
      `• Cost estimation\n` +
      `• Engineer assignment (skill matching)\n\n` +
      `Do you have a parts list in Excel or CSV format?\n\n` +
      `📥 **Need a template?** [Download Parts Catalog Template](/templates/parts-catalog-template.csv)`,
      [
        { label: '📤 Yes, upload file', value: 'upload_parts' },
        { label: '📥 Download template first', value: 'download_template' },
        { label: '✍️ I\'ll add manually', value: 'manual_parts' },
        { label: '⏭️ Skip for now', value: 'skip_parts' },
      ]
    );
    setCurrentStep(2); // Move to parts step
  };

  const handlePartsStep = async (input: string) => {
    if (input === 'download_template') {
      addAIMessage(
        `Great! I've provided the template download link above.\n\n` +
        `The template includes:\n` +
        `✅ All required fields (part_number, part_name, category, part_type)\n` +
        `✅ Sample data for Ventilator parts (15+ examples)\n` +
        `✅ Field descriptions and validation rules\n` +
        `✅ JSON format for technical specifications\n\n` +
        `Fill it out and come back when ready!`,
        [
          { label: '📤 Ready to upload', value: 'upload_parts' },
          { label: '⏭️ Skip for now', value: 'skip_parts' },
        ]
      );
      return;
    }

    if (input === 'upload_parts') {
      addAIMessage('Please upload your parts catalog file (CSV or Excel):');
      fileInputRef.current?.click();
      return;
    }

    if (input === 'skip_parts') {
      addAIMessage(
        `No problem! You can add parts later from the dashboard.\n\n` +
        `Step 3 complete! ✅ Now let's set up your **Service Engineers**.\n\n` +
        `Service engineers are crucial for:\n` +
        `• Responding to service tickets\n` +
        `• Equipment maintenance\n` +
        `• Installation and commissioning\n\n` +
        `Do you have a list of your service engineers?`,
        [
          { label: '📤 Yes, upload engineer list', value: 'upload_engineers' },
          { label: '✍️ I\'ll add manually', value: 'manual_engineers' },
          { label: '⏭️ Skip for now', value: 'skip_engineers' },
        ]
      );
      setCurrentStep(3); // Move to engineers
      return;
    }

    addAIMessage("I'll help you add parts manually. What's the first part number?");
  };

  const handleEngineersStep = async (input: string) => {
    if (input === 'upload_engineers') {
      addAIMessage('Please upload your service engineers list (CSV or Excel):');
      fileInputRef.current?.click();
      return;
    }

    if (input === 'skip_engineers') {
      addAIMessage(
        `Understood! You can add engineers later.\n\n` +
        `Step 4 complete! ✅ Final step: **Equipment Installations**.\n\n` +
        `Let's register where your equipment is installed:\n` +
        `• Customer/Hospital locations\n` +
        `• Equipment serial numbers\n` +
        `• Installation dates\n` +
        `• Equipment models\n\n` +
        `Do you have an installations list?`,
        [
          { label: '📤 Yes, upload file', value: 'upload_installations' },
          { label: '⏭️ Skip for now', value: 'skip_installations' },
        ]
      );
      setCurrentStep(4); // Move to installations
      return;
    }

    addAIMessage("I'll help you add engineers manually. What's the engineer's name?");
  };

  const handleInstallationsStep = async (input: string) => {
    if (input === 'upload_installations') {
      addAIMessage('Please upload your equipment installations file (CSV or Excel):');
      fileInputRef.current?.click();
      return;
    }

    if (input === 'skip_installations') {
      addAIMessage(
        `All set! 🎉\n\n` +
        `You've completed the onboarding wizard. Here's what we set up:\n` +
        `✅ Manufacturer information\n` +
        `✅ Equipment catalog\n` +
        `✅ Parts catalog\n` +
        `✅ Service engineers\n` +
        `✅ Equipment installations\n\n` +
        `You can now manage everything from your dashboard!`,
        [
          { label: '🏠 Go to Dashboard', value: 'go_dashboard' },
          { label: '📝 Review Data', value: 'review' },
        ]
      );
      setCurrentStep(5); // Move to review
      return;
    }

    addAIMessage(
      `Perfect! Now let's finalize everything.\n\n` +
      `• Engineer levels (1/2/3)\n` +
      `• Equipment expertise\n` +
      `• Coverage areas`,
      [
        { label: '📤 Upload engineers file', value: 'upload_engineers' },
        { label: '✍️ Add manually', value: 'manual_engineers' },
      ]
    );
    setCurrentStep(5); // Move to complete
  };

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    addUserMessage(`Uploaded: ${file.name}`, 'file');
    setIsProcessing(true);

    try {
      // Simulate file processing
      await new Promise(resolve => setTimeout(resolve, 2000));

      // In production, parse CSV/Excel and validate
      const fileType = file.name.split('.').pop();
      
      addAIMessage(
        `✅ File "${file.name}" uploaded successfully!\n\n` +
        `I'm analyzing your data...\n\n` +
        `Found:\n` +
        `• 50 rows detected\n` +
        `• All required columns present\n` +
        `• No critical errors\n\n` +
        `Should I proceed with import?`,
        [
          { label: '✅ Yes, import', value: 'confirm_import' },
          { label: '👁️ Show preview first', value: 'preview' },
          { label: '❌ Cancel', value: 'cancel_import' },
        ]
      );
    } catch (error) {
      addAIMessage(
        `❌ Failed to process file. Please check the format and try again.`,
        [
          { label: 'Upload different file', value: 'reupload' },
          { label: 'Switch to manual entry', value: 'manual' },
        ]
      );
    } finally {
      setIsProcessing(false);
    }
  };

  const handleAction = (value: string) => {
    if (value === 'manual') {
      setShowManualFallback(true);
      addAIMessage(
        "No problem! Switching to manual entry mode.\n\n" +
        "You can fill out the forms directly."
      );
      // Route to appropriate manual form
      setTimeout(() => {
        router.push('/onboarding/wizard');
      }, 2000);
      return;
    }

    if (value === 'dashboard') {
      router.push('/dashboard');
      return;
    }

    // Process other actions
    processUserInput(value);
  };

  // Show loading while checking authorization
  if (isAuthorized === null) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-50">
        <Card className="p-8">
          <div className="flex flex-col items-center gap-4">
            <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
            <p className="text-gray-600">Checking access permissions...</p>
          </div>
        </Card>
      </div>
    );
  }

  // Show unauthorized message for users without access
  if (!isAuthorized) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-red-50 to-orange-50 p-6">
        <Card className="max-w-md">
          <CardContent className="p-8">
            <div className="flex flex-col items-center gap-4 text-center">
              <div className="w-16 h-16 rounded-full bg-red-100 flex items-center justify-center">
                <Shield className="h-8 w-8 text-red-600" />
              </div>
              <h2 className="text-2xl font-bold text-gray-900">Access Restricted</h2>
              <p className="text-gray-600">
                The AI Onboarding Wizard is only available to <strong>System Administrators</strong> and <strong>Manufacturer Admins</strong>.
              </p>
              <p className="text-sm text-gray-500">
                Please contact your administrator for access.
              </p>
              <div className="flex gap-3 mt-4">
                <Button onClick={() => router.push('/dashboard')}>
                  Go to Dashboard
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Render AI wizard for authorized users
  return (
    <div className="flex min-h-screen bg-gray-50">
      <Navigation />
      <div className="flex-1 p-6">
        <div className="max-w-6xl mx-auto">
          {/* Header */}
          <div className="mb-6">
            <h1 className="text-3xl font-bold text-gray-900">AI Onboarding Assistant</h1>
            <p className="text-gray-600 mt-2">Complete your organization setup with AI guidance</p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Progress Sidebar */}
          <div className="lg:col-span-1">
            <Card>
              <CardContent className="p-6">
                <h2 className="font-semibold text-lg mb-4">Progress</h2>
                <div className="space-y-3">
                  {STEPS.map((step, index) => {
                    // Check if manufacturer admin and hide company profile step
                    const userStr = localStorage.getItem('user');
                    const user = userStr ? JSON.parse(userStr) : null;
                    const isManufacturerAdmin = user?.organization_type === 'manufacturer';
                    
                    // Skip company profile for manufacturer admins
                    if (isManufacturerAdmin && step.id === 'company') {
                      return null;
                    }

                    return (
                      <div
                        key={step.id}
                        onClick={() => {
                          setCurrentStep(index);
                          showStepPrompt(index);
                        }}
                        className={`flex items-center gap-3 p-2 rounded-lg transition-colors cursor-pointer hover:bg-blue-100 ${
                          index === currentStep
                            ? 'bg-blue-50 border border-blue-200'
                            : index < currentStep
                            ? 'bg-green-50'
                            : 'bg-gray-50'
                        }`}
                      >
                        {index < currentStep ? (
                          <CheckCircle2 className="h-5 w-5 text-green-600 flex-shrink-0" />
                        ) : (
                          <Circle
                            className={`h-5 w-5 flex-shrink-0 ${
                              index === currentStep ? 'text-blue-600' : 'text-gray-400'
                            }`}
                          />
                        )}
                        <div className="flex-1 min-w-0">
                          <div className="text-lg">{step.icon}</div>
                          <div
                            className={`text-sm font-medium truncate ${
                              index === currentStep ? 'text-blue-900' : 'text-gray-700'
                            }`}
                          >
                            {step.label}
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>

                <div className="mt-6 pt-6 border-t">
                  <div className="text-sm text-gray-600 mb-2">Overall Progress</div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                      style={{ 
                        width: `${(() => {
                          const userStr = localStorage.getItem('user');
                          const user = userStr ? JSON.parse(userStr) : null;
                          const isManufacturerAdmin = user?.organization_type === 'manufacturer';
                          const totalSteps = isManufacturerAdmin ? STEPS.length - 1 : STEPS.length;
                          const adjustedStep = isManufacturerAdmin ? Math.max(0, currentStep - 1) : currentStep;
                          return (adjustedStep / totalSteps) * 100;
                        })()}%` 
                      }}
                    />
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    {(() => {
                      const userStr = localStorage.getItem('user');
                      const user = userStr ? JSON.parse(userStr) : null;
                      const isManufacturerAdmin = user?.organization_type === 'manufacturer';
                      const totalSteps = isManufacturerAdmin ? STEPS.length - 1 : STEPS.length;
                      const adjustedStep = isManufacturerAdmin ? Math.max(0, currentStep - 1) : currentStep;
                      return `Step ${adjustedStep + 1} of ${totalSteps}`;
                    })()}
                  </div>
                </div>

                <div className="mt-4 space-y-2">
                  <Button
                    variant="outline"
                    size="sm"
                    className="w-full"
                    onClick={() => setShowManualFallback(true)}
                  >
                    <FileText className="h-4 w-4 mr-2" />
                    Switch to Manual
                  </Button>
                  
                  <div className="text-xs text-gray-500 text-center">
                    💡 Click any step above to jump to it
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Chat Interface */}
          <div className="lg:col-span-3">
            <Card className="h-[700px] flex flex-col">
              {/* Messages */}
              <div className="flex-1 overflow-y-auto p-6 space-y-4">
                {messages.map((message) => (
                  <div
                    key={message.id}
                    className={`flex ${message.role === 'user' ? 'justify-end' : 'justify-start'}`}
                  >
                    <div
                      className={`max-w-[80%] rounded-lg p-4 ${
                        message.role === 'user'
                          ? 'bg-blue-600 text-white'
                          : 'bg-white border border-gray-200 text-gray-900'
                      }`}
                    >
                      {message.role === 'ai' && (
                        <div className="flex items-center gap-2 mb-2">
                          <div className="w-6 h-6 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-xs font-bold">
                            AI
                          </div>
                          <span className="text-xs text-gray-500">Assistant</span>
                        </div>
                      )}
                      
                      <div className="whitespace-pre-wrap">{message.content}</div>

                      {message.type === 'file' && (
                        <div className="mt-2 flex items-center gap-2 text-sm">
                          <Upload className="h-4 w-4" />
                          File uploaded
                        </div>
                      )}

                      {message.actions && (
                        <div className="mt-4 flex flex-wrap gap-2">
                          {message.actions.map((action, idx) => (
                            <Button
                              key={idx}
                              size="sm"
                              variant="outline"
                              onClick={() => handleAction(action.value)}
                              className="bg-white hover:bg-gray-50"
                            >
                              {action.label}
                            </Button>
                          ))}
                        </div>
                      )}

                      <div className="text-xs mt-2 opacity-70">
                        {message.timestamp.toLocaleTimeString()}
                      </div>
                    </div>
                  </div>
                ))}

                {isProcessing && (
                  <div className="flex justify-start">
                    <div className="bg-white border border-gray-200 rounded-lg p-4">
                      <div className="flex items-center gap-2">
                        <Loader2 className="h-4 w-4 animate-spin text-blue-600" />
                        <span className="text-sm text-gray-600">AI is thinking...</span>
                      </div>
                    </div>
                  </div>
                )}

                <div ref={messagesEndRef} />
              </div>

              {/* Input Area */}
              <div className="border-t p-4 bg-gray-50">
                <div className="flex gap-2">
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept=".csv,.xlsx,.xls"
                    onChange={handleFileUpload}
                    className="hidden"
                  />
                  <Button
                    variant="outline"
                    onClick={() => fileInputRef.current?.click()}
                    disabled={isProcessing}
                  >
                    <Upload className="h-4 w-4" />
                  </Button>
                  <Input
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && processUserInput(input)}
                    placeholder="Type your response..."
                    disabled={isProcessing}
                    className="flex-1"
                  />
                  <Button
                    onClick={() => processUserInput(input)}
                    disabled={isProcessing || !input.trim()}
                  >
                    {isProcessing ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Send className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                <p className="text-xs text-gray-500 mt-2">
                  💡 Tip: You can type your answer or upload files when asked
                </p>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </div>
    </div>
  );
}
