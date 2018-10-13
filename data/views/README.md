# Application View Folder
This folder houses the full page version of each modals views. In other words it includes the Dochead and docfooter as well as the sidebar navgation, topbar with breadcrumbs.

# Widget View Folder
This folder houses the partial views for each of the modals. They must be included in a modal application view. These are the individual building blocks for a model.

### Model Specific
* dash - Shows a simple dashboard for the current model (stats and graphs)
* find - Profiles filters and seraching for the model
* form - Allows for adding new models or editing existing models
* full - Provides a standard (across all models) full view (dash, find, menu, list)
* item - Displays a single model (with proper access provides delete and edit options)
* list - Display all models in a responsive table. (with proper access provides view, delete and edit options for each model in the list)
* menu - provides a common menu for the model
- modal - Popup Modal dialog with full widget as the body

### Misc Common Widgets
* calendar - Provides a calendar
* eventcalendar - provides a user/group calendars with the ability to add/remove events
* googleanalytics - if enabled in config, will embed google analytics for the webapp
* carousel - used to display DB defined carousels/onboardings

### Model Widgets
* TODO

# Root Folder
Common partial templates that are used throughout the entire webapp.
* dochead - The top part (before content) of the HTML5 document, includes topnavbar, googleanaylticvs and sidebarmenu
* docfoot - The bottom part (after content) for the HTML5 document
* sidebarmenu - Provides the menu you see in the left side that can open and close
* topnavbar - provides the top bar with notifications, sidebar open/close, user account menu as well as the breadcrumbs
