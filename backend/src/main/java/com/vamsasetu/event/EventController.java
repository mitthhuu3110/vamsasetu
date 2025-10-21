package com.vamsasetu.event;

import com.vamsasetu.common.ApiResponse;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/events")
public class EventController {
    @GetMapping()
    public ResponseEntity<ApiResponse<?>> getAll() {
        // TODO: fetch events
        return ResponseEntity.ok(new ApiResponse<>(true, "Events", null));
    }
    @PostMapping()
    public ResponseEntity<ApiResponse<?>> add(@RequestBody Object event) {
        // TODO: add event
        return ResponseEntity.ok(new ApiResponse<>(true, "Added", null));
    }
    // ...
}
