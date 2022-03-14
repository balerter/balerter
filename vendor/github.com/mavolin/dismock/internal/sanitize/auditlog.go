package sanitize

import (
	"github.com/diamondburned/arikawa/discord"
)

// AuditLog sanitizes an discord.AuditLog.
//
// This function will sanitize AuditLog.Webhook.ID, AuditLog.Webhook.User.ID,
// AuditLog.Users.ID, AuditLog.Entries.ID, AuditLog.Entries.UserID,
// AuditLog.Integrations.ID, AuditLog.Integrations.RoleID and
// AuditLog.Integrations.User.ID.
func AuditLog(al discord.AuditLog) discord.AuditLog {
	for i, w := range al.Webhooks {
		al.Webhooks[i] = Webhook(w, 1, 1, 1)
	}

	for i, u := range al.Users {
		al.Users[i] = User(u, 1)
	}

	for i, e := range al.Entries {
		al.Entries[i] = AuditLogEntry(e, 1, 1)
	}

	for i, integration := range al.Integrations {
		al.Integrations[i] = Integration(integration, 1, 1, 1)
	}

	return al
}

// AuditLogEntry sanitizes an discord.AuditLogEntry.
//
// This function will sanitize AuditLogEntry.ID and AuditLogEntry.UserID.
func AuditLogEntry(e discord.AuditLogEntry, id discord.AuditLogEntryID, userID discord.UserID) discord.AuditLogEntry {
	if e.ID == 0 {
		e.ID = id
	}

	if e.UserID == 0 {
		e.UserID = userID
	}

	return e
}
