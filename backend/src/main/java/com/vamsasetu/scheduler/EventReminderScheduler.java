package com.vamsasetu.scheduler;

import com.vamsasetu.model.Event;
import com.vamsasetu.model.EventAttendee;
import com.vamsasetu.repository.EventRepository;
import com.vamsasetu.service.NotificationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;

@Component
public class EventReminderScheduler {

    @Autowired
    private EventRepository eventRepository;

    @Autowired
    private NotificationService notificationService;

    @Scheduled(cron = "0 0 9 * * ?") // Run daily at 9 AM
    public void sendDailyEventReminders() {
        LocalDate today = LocalDate.now();
        LocalDate tomorrow = today.plusDays(1);
        
        // Find events happening today or tomorrow that have reminders enabled
        List<Event> upcomingEvents = eventRepository.findEventsWithRemindersInDateRange(
            null, // This should be filtered by user in a real implementation
            today,
            tomorrow.plusDays(1)
        );
        
        for (Event event : upcomingEvents) {
            if (shouldSendReminder(event)) {
                sendRemindersForEvent(event);
            }
        }
    }

    private boolean shouldSendReminder(Event event) {
        if (!event.getReminderEnabled()) {
            return false;
        }
        
        // Check if reminder should be sent based on advance time
        LocalDateTime reminderTime = event.getEventDate().atTime(
            event.getEventTime() != null ? event.getEventTime() : event.getEventDate().atStartOfDay().toLocalTime()
        ).minusHours(event.getReminderAdvanceTime());
        
        LocalDateTime now = LocalDateTime.now();
        return now.isAfter(reminderTime) && now.isBefore(reminderTime.plusHours(1));
    }

    private void sendRemindersForEvent(Event event) {
        for (EventAttendee attendee : event.getAttendees()) {
            if (attendee.getMember() != null) {
                notificationService.sendEventReminder(event, attendee.getMember());
            }
        }
    }
}
