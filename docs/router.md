- For APIs with different permission levels for guest and staff, add `/admin` before the permissions. If there's only one permission level, you don't need to add anything for both guest and staff.

- For APIs that use a basic list to retrieve a dropdown, add `?view=basic`. This returns a list without page sections, containing only the ID and name.