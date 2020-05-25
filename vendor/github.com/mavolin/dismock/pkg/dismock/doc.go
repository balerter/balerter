// Package dismock creates mocks for the Discord API.
// The names of the mocks correspond to arikawa's API wrapper names, put as
// this are http mocks, any discord library can be mocked.
//
// Field Sanitation
//
// As you might have noticed, some of the MockX methods have footers like:
//		This method will sanitize emoji.ID and emoji.User.ID.
// This means that all fields mentioned in the comment, typically all
// Snowflakes that are not omittable, will be set to 1, or a
// corresponding value passed in as a parameter, if their value is v <= 0.
// This is necessary, as arikawa serializes all Snowflakes that are s <= 0 to
// JSON null, as they are seen as invalid.
//
// However, this shouldn't impose much of a problem as 0 or smaller isn't a
// valid Snowflake anyway, and all values above 0 will not be sanitized.
//
// Mocking Requests for Metadata
//
// In order to mock requests for an entity's meta data, such as it's icon, it
// is required that those requests are made with Mocker.Client, so that
// the request are correctly redirected.
//
//
// Important notes
//
// BUG(mavolin): Due to an inconvenient behavior of json.Unmarshal, where on
// JSON null the go representation's MarshalJSON method doesn't get called,
// there is no way to differentiate between option.NullX and that option type
// omitted.
package dismock
