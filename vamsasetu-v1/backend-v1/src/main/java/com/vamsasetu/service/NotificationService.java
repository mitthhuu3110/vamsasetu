package com.vamsasetu.service;

import com.vamsasetu.model.Event;
import com.vamsasetu.model.FamilyMember;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.stereotype.Service;

import java.time.format.DateTimeFormatter;

@Service
public class NotificationService {

    @Autowired
    private JavaMailSender mailSender;

    public void sendEventReminder(Event event, FamilyMember member) {
        if (event.getReminderEmail() && member.getEmail() != null) {
            sendEmailReminder(event, member);
        }
        
        if (event.getReminderSms() && member.getPhone() != null) {
            sendSmsReminder(event, member);
        }
        
        if (event.getReminderWhatsapp() && member.getPhone() != null) {
            sendWhatsappReminder(event, member);
        }
    }

    private void sendEmailReminder(Event event, FamilyMember member) {
        try {
            SimpleMailMessage message = new SimpleMailMessage();
            message.setTo(member.getEmail());
            message.setSubject("Reminder: " + event.getTitle());
            
            String body = buildEmailBody(event, member);
            message.setText(body);
            
            mailSender.send(message);
            System.out.println("Email reminder sent to: " + member.getEmail());
        } catch (Exception e) {
            System.err.println("Failed to send email reminder: " + e.getMessage());
        }
    }

    private void sendSmsReminder(Event event, FamilyMember member) {
        // TODO: Implement SMS using Twilio
        System.out.println("SMS reminder would be sent to: " + member.getPhone() + " for event: " + event.getTitle());
    }

    private void sendWhatsappReminder(Event event, FamilyMember member) {
        // TODO: Implement WhatsApp using Twilio or Meta Graph API
        System.out.println("WhatsApp reminder would be sent to: " + member.getPhone() + " for event: " + event.getTitle());
    }

    private String buildEmailBody(Event event, FamilyMember member) {
        StringBuilder body = new StringBuilder();
        body.append("Dear ").append(member.getFirstName()).append(",\n\n");
        body.append("This is a reminder for the upcoming event:\n\n");
        body.append("Event: ").append(event.getTitle()).append("\n");
        body.append("Date: ").append(event.getEventDate().format(DateTimeFormatter.ofPattern("MMMM dd, yyyy"))).append("\n");
        
        if (event.getEventTime() != null) {
            body.append("Time: ").append(event.getEventTime().format(DateTimeFormatter.ofPattern("hh:mm a"))).append("\n");
        }
        
        if (event.getLocation() != null) {
            body.append("Location: ").append(event.getLocation()).append("\n");
        }
        
        if (event.getDescription() != null) {
            body.append("\nDescription: ").append(event.getDescription()).append("\n");
        }
        
        body.append("\nBest regards,\n");
        body.append("VamsaSetu Team");
        
        return body.toString();
    }
}
