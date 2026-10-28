- Add new tenant:
  + The .env.prod file must contain the following variables: TA_DOMAIN, TA_SCHEMA, SD_NAME=tenant

- Initialize administrator for the tenant:
  + The .env.prod file must contain the following variables: AD_SCHEMA, AD_USERNAME, AD_EMAIL, AD_PASSWORD

- Pricing calendar:
  + `SD_NAME=tenant` will automatically seed pricing calendar for 2 years (current year and next year) for the newly created tenant only.
