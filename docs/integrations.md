# Integrations

## Home Assistant

Track anything from Home Assistant automations — door opens, workouts, coffee brews, whatever.

### Prerequisites

- A running Daytrack instance (self-hosted)
- An [API key](../README.md#api-keys) from the web UI

### RESTful Command (simplest)

```yaml
action: rest_command.daytrack_track
```

Add to `configuration.yaml`:

```yaml
rest_command:
  daytrack_track:
    url: "https://daytrack.example.com/v1/events/{{ username }}/{{ track }}?key={{ api_key }}&quantity={{ quantity }}"
    method: POST
```

Then in any automation:

```yaml
alias: "Track morning coffee"
trigger:
  - platform: state
    entity_id: binary_sensor.coffee_machine
    to: "on"
action:
  - action: rest_command.daytrack_track
    data:
      username: "your_username"
      track: "coffee"
      api_key: "YOUR_API_KEY"
      quantity: 1
```

### Direct curl (no REST command setup)

```yaml
alias: "Log workout on Apple Watch notification"
trigger:
  - platform: event
    event_type: mobile_app_notification_action
    event_data:
      action: "WORKOUT_DONE"
action:
  - action: shell_command.curl_track
    data:
      track: "workout"
```

Add to `configuration.yaml`:

```yaml
shell_command:
  curl_track: "curl -X POST 'https://daytrack.example.com/v1/events/{{ username }}/{{ track }}?key={{ api_key }}&quantity=1'"
```

### Tracking with custom timestamps

Useful for backfilling or tracking past events:

```yaml
action: rest_command.daytrack_track
data:
  username: "your_username"
  track: "reading"
  api_key: "YOUR_API_KEY"
  quantity: 15
  # Override with created_at={{ now().isoformat() }} in the URL template
```

The `created_at` parameter accepts RFC 3339 format: `&created_at=2026-06-03T10:00:00Z`

### Public tracks (no API key needed)

If you set a track's visibility to `public_rw` (public read+write), you can omit the API key:

```yaml
action: rest_command.daytrack_track
data:
  username: "your_username"
  track: "my_public_track"
  api_key: ""  # not needed
  quantity: 1
```

---

## Flic (physical button)

Track anything with a button press — no phone needed.

### Prerequisites

- Flic hub with internet access
- An [API key](../README.md#api-keys) from the web UI

### Flic app setup

1. Open the Flic app
2. Select your button
3. Tap **"Add action"** → **"Internet Request"**
4. Configure:

| Field | Value |
|---|---|
| **URL** | `https://daytrack.example.com/v1/events/{username}/{track_name}?key={api_key}&quantity=1` |
| **Method** | `POST` |
| **Headers** | (leave empty) |
| **Body** | (leave empty) |

Replace:
- `{username}` — your Daytrack username
- `{track_name}` — the track name (e.g. `water`, `standup`, `coffee`)
- `{api_key}` — your API key

### Multiple tracks with one button (Flic long press)

Assign different tracks to different press types:

| Press | URL |
|---|---|
| **Click** | `https://.../coffee?key=xxx&quantity=1` |
| **Double click** | `https://.../water?key=xxx&quantity=1` |
| **Hold** | `https://.../standup?key=xxx&quantity=1` |

### Public tracks (no API key)

For public-writable tracks, omit the `key` parameter:

```
https://daytrack.example.com/v1/events/{username}/{track_name}?quantity=1
```

---

## Custom (any HTTP client)

```bash
# Track a thing
curl -X POST "https://daytrack.example.com/v1/events/john/coffee?key=abc123&quantity=1"

# Track with timestamp
curl -X POST "https://daytrack.example.com/v1/events/john/coffee?key=abc123&quantity=1&created_at=2026-06-03T10:00:00Z"

# List recent events
curl "https://daytrack.example.com/v1/events/john/coffee?key=abc123&list_by=day"
```

Works from:
- **Shortcuts (iOS)** — "Get contents of URL" action, POST
- **Tasker (Android)** — HTTP Post action
- **IFTTT** — Webhooks → Maker channel
- **n8n / Node-RED** — HTTP Request node
- **Any curl-capable device** — ESP32, Raspberry Pi, etc.
