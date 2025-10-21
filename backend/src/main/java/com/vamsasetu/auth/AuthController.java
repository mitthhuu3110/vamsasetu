package com.vamsasetu.auth;

import com.vamsasetu.common.ApiResponse;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/auth")
public class AuthController {
    @PostMapping("/register")
    public ResponseEntity<ApiResponse<?>> register(@RequestBody Object registerRequest) {
        // TODO: registration logic
        return ResponseEntity.ok(new ApiResponse<>(true, "Registered", null));
    }

    @PostMapping("/login")
    public ResponseEntity<ApiResponse<?>> login(@RequestBody Object loginRequest) {
        // TODO: login logic
        return ResponseEntity.ok(new ApiResponse<>(true, "Logged in", null));
    }

    @GetMapping("/profile")
    public ResponseEntity<ApiResponse<?>> profile() {
        // TODO: profile logic
        return ResponseEntity.ok(new ApiResponse<>(true, "Profile", null));
    }
}
