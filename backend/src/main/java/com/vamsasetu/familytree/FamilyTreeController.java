package com.vamsasetu.familytree;

import com.vamsasetu.common.ApiResponse;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/family")
public class FamilyTreeController {
    @GetMapping("/members")
    public ResponseEntity<ApiResponse<?>> getAllMembers() {
        // TODO: fetch members
        return ResponseEntity.ok(new ApiResponse<>(true, "Fetched", null));
    }
    @PostMapping("/members")
    public ResponseEntity<ApiResponse<?>> addMember(@RequestBody Object member) {
        // TODO: add member
        return ResponseEntity.ok(new ApiResponse<>(true, "Added", null));
    }
    @PostMapping("/relations")
    public ResponseEntity<ApiResponse<?>> addRelation(@RequestBody Object relation) {
        // TODO: add relation
        return ResponseEntity.ok(new ApiResponse<>(true, "Relation Added", null));
    }
    // More endpoints ...
}
