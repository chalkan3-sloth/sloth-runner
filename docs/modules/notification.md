# Notification Module

Send notifications to various channels (Slack, Discord, Email, etc).

## Overview

The notification module provides a unified interface for sending alerts and notifications to different platforms.

## Functions

### `notification.slack(webhook, message)`

Send a message to Slack.

```lua
notification.slack(
  "https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
  "Deployment completed successfully!"
)
```

### `notification.discord(webhook, message)`

Send a message to Discord.

```lua
notification.discord(
  "https://discord.com/api/webhooks/YOUR/WEBHOOK",
  "Build failed!"
)
```

### `notification.email(config)`

Send an email notification.

```lua
notification.email({
  to = "team@example.com",
  subject = "Deployment Alert",
  body = "Deployment to production completed"
})
```

## See Also

- [Notifications Documentation](../en/modules/notifications.md)
- [Alert Configuration](../en/advanced-features.md)
