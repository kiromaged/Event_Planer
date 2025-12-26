CREATE DATABASE IF NOT EXISTS `EventPlanner`
DEFAULT CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

USE `EventPlanner`;

-------------------------------------------------------
-- USERS TABLE
-------------------------------------------------------
CREATE TABLE IF NOT EXISTS `users` (
    `user_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password_hash` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `ux_users_email` (`email`),
    KEY `ix_users_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-------------------------------------------------------
-- EVENTS TABLE
-------------------------------------------------------
CREATE TABLE IF NOT EXISTS `events` (
    `event_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `location` VARCHAR(255) NOT NULL,
    `event_date` DATE NOT NULL,
    `event_time` TIME NOT NULL,
    `created_by` INT UNSIGNED NOT NULL, -- the main organizer
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (`event_id`),

    KEY `ix_events_title` (`title`),
    KEY `ix_events_date` (`event_date`),
    KEY `ix_events_creator` (`created_by`),

    CONSTRAINT `fk_events_creator`
        FOREIGN KEY (`created_by`) REFERENCES `users` (`user_id`)
        ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-------------------------------------------------------
-- EVENT ATTENDEES TABLE
-- Stores invitations + roles per event
-------------------------------------------------------
CREATE TABLE IF NOT EXISTS `event_attendees` (
    `event_id` INT UNSIGNED NOT NULL,
    `user_id` INT UNSIGNED NOT NULL,

    -- organizer, attendee
    `role` ENUM('organizer','attendee') NOT NULL DEFAULT 'attendee',

    -- attendance status
    `status` ENUM('going','maybe','not_going','pending') 
        NOT NULL DEFAULT 'pending',

    `invited_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (`event_id`, `user_id`),

    KEY `ix_attendees_status` (`status`),
    KEY `ix_attendees_user` (`user_id`),
    KEY `ix_attendees_role` (`role`),
    KEY `ix_attendees_user_role` (`user_id`, `role`),

    CONSTRAINT `fk_event_attendees_event`
        FOREIGN KEY (`event_id`) REFERENCES `events` (`event_id`)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT `fk_event_attendees_user`
        FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
        ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-------------------------------------------------------
-- TASKS TABLE
-- Tasks associated with events
-------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tasks` (
    `task_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `event_id` INT UNSIGNED NOT NULL,
    `description` TEXT NOT NULL,
    `assigned_to` INT UNSIGNED NULL, -- NULL if unassigned
    `status` ENUM('pending','in_progress','completed','cancelled') 
        NOT NULL DEFAULT 'pending',
    `due_date` DATE NULL,
    `created_by` INT UNSIGNED NOT NULL, -- who created the task
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`task_id`),

    KEY `ix_tasks_event` (`event_id`),
    KEY `ix_tasks_assigned` (`assigned_to`),
    KEY `ix_tasks_status` (`status`),
    KEY `ix_tasks_due_date` (`due_date`),
    KEY `ix_tasks_creator` (`created_by`),

    CONSTRAINT `fk_tasks_event`
        FOREIGN KEY (`event_id`) REFERENCES `events` (`event_id`)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT `fk_tasks_assigned_user`
        FOREIGN KEY (`assigned_to`) REFERENCES `users` (`user_id`)
        ON DELETE SET NULL ON UPDATE CASCADE,

    CONSTRAINT `fk_tasks_creator`
        FOREIGN KEY (`created_by`) REFERENCES `users` (`user_id`)
        ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-------------------------------------------------------
-- SEARCH SUPPORT: Optional Keywords Table (Optional)
-- Helps with advanced filtering & indexing for search
-------------------------------------------------------
CREATE TABLE IF NOT EXISTS `event_keywords` (
    `keyword_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `event_id` INT UNSIGNED NOT NULL,
    `keyword` VARCHAR(100) NOT NULL,

    PRIMARY KEY (`keyword_id`),
    KEY `ix_keyword_event` (`event_id`),
    FULLTEXT KEY `ft_keyword_text` (`keyword`),

    CONSTRAINT `fk_event_keywords_event`
        FOREIGN KEY (`event_id`) REFERENCES `events` (`event_id`)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-------------------------------------------------------
-- INDEXES FOR SEARCH (Full Text)
-------------------------------------------------------
ALTER TABLE `events`
    ADD FULLTEXT KEY `ft_events_title_desc` (`title`, `description`);

ALTER TABLE `tasks`
    ADD FULLTEXT KEY `ft_tasks_description` (`description`);
