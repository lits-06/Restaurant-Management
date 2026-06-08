# User Service HTTP API

Current status: no HTTP endpoints are exposed via API Gateway for User Service yet.

## Front-end note
- Do not call User Service directly over HTTP at this time.
- User/profile/role flows should be considered unavailable in HTTP layer until gateway routes are added.

## Planned scope (from project plan)
- Create user/staff
- Update user/staff info
- Delete user/staff
- Role management (RBAC)

When HTTP routes are added, this file should be updated with:
- Endpoint list
- Request body/input schemas
- Response/output schemas
- Error handling samples
