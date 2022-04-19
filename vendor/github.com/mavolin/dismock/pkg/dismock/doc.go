// Package dismock creates mocks for the Discord API.
// The names of the mocks correspond to arikawa's API wrapper names, but as
// this are http mocks, any discord library can be mocked.
//
// Field Sanitation
//
// As you might have noticed, some of the MockX methods have footers like:
//		This method will sanitize Emoji.ID and Emoji.User.ID.
// This means that all fields mentioned in the comment will be set to 1, or, if
// available, a value passed in as a parameter, if their value is v <= 0.
// This is necessary, as arikawa serializes all Snowflakes that are s <= 0 to
// JSON null, as they are seen as invalid.
//
// However, this shouldn't impose much of a problem as a Snowflake with the
// value 0 or smaller isn't valid anyway, and all valid values 0 will not be
// sanitized.
//
// Mocking Requests for Metadata
//
// Besides the regular API calls, dismock also features mocks for fetching
// an entities meta data, e.g. an icon or a splash.
// In order to mock requests for an entity's meta data, you need to make sure
// that those requests are made with Mocker.Client, so that the requests are
// correctly redirected to the mock server.
//
// Mocking Errors
//
// Not always do we expect an API call to succeed.
// To send a discord error, use the Mocker.Error method.
//
//
// Important Notes
//
// BUG(mavolin): Due to an inconvenient behavior of json.Unmarshal, where on
// JSON null the the MarshalJSON method doesn't get called, there is no way to
// differentiate between option.NullX and omitted.
package dismock
