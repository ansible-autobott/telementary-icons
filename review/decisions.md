<!--
review/decisions.md — written by the telementary overlay review page
(make review, then http://localhost:3033/review), which rewrites the Approved
and Rejected categories below on every decision. It is also a todo-app list:
todo review/decisions.md

For the implementer: each task proposes publishing elementary art (paths are
relative to the elementary icon theme root) at the listed theme buckets. A line
"24, 22 ← actions/24/x.svg" is one overlay.yaml link: copy that file into an
overlays/<category>/<set>/ folder and publish it at sizes [24, 22] (the first
is the real copy). "(→ …)" marks an elementary symlink: copy the file it
resolves to. Symbolic art must be recolored to the ColorScheme-Text pattern
(overlays/README.md). "also changes" lists Tela aliases that follow the
replaced file. A Comment overrides the proposal. Approved: implement it.
Rejected: do not overlay it, but still act on its comment if it asks for
something else. Mark a task [x] once it is done.
-->

# Elementary overlay review

## Approved

- [x] `document-save` — elementary overlay
  names: document-save, document-save-symbolic · category: actions
  16 ← actions/16/document-save.svg
  24, 22 ← actions/24/document-save.svg
  32 ← actions/32/document-save.svg
  48 ← actions/48/document-save.svg
  document-save-symbolic: symbolic ← actions/symbolic/document-save-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: browser-download (16, 22, 24), epiphany-download (16, 22, 24), filesave (16, 22, 24), gtk-floppy (16, 22, 24), gtk-save (16, 22, 24), modified (16, 22, 24), stock_save (16, 22, 24), system-save-session (16, 22, 24), xsi-document-save-symbolic (symbolic)
- [x] `document-new` — elementary overlay
  names: document-new, document-new-symbolic · category: actions
  16 ← actions/16/document-new.svg
  24, 22 ← actions/24/document-new.svg
  32 ← actions/32/document-new.svg
  48 ← actions/48/document-new.svg
  document-new-symbolic: symbolic ← actions/symbolic/document-new-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: add-files-to-archive (16, 22, 24), archive-insert (16, 22, 24), document-single (16), filenew (16, 22, 24), stock_new-text (16, 22, 24), vm-create (16, 22, 24), xsi-document-new-symbolic (symbolic)
- [x] `document-edit` — elementary overlay
  names: document-edit, document-edit-symbolic · category: actions
  16 ← actions/16/document-edit.svg
  24, 22 ← actions/24/document-edit.svg
  48 ← actions/48/document-edit.svg
  document-edit-symbolic: symbolic ← actions/symbolic/document-edit-symbolic.svg (recolor to ColorScheme-Text)
  also changes: cell_edit (16, 22, 24), document-sign (16, 22, 24), edit-entry (16, 22, 24), edit-map (16, 22, 24), editimage (16, 22, 24), editor (16, 22, 24), edittext (16, 22, 24), gnumeric-comment-edit (16, 22, 24), gnumeric-link-edit (16, 22, 24), media-track-edit-amarok (16, 22, 24), open-for-editing (16, 22, 24), xsi-document-edit-symbolic (symbolic)
- [x] `document-open-folder` — elementary overlay from `folder-open`
  names: document-open-folder · category: actions
  16 ← places/16/folder-open.svg
  24, 22 ← places/24/folder-open.svg
  32 ← places/32/folder-open.svg
  48 ← places/48/folder-open.svg
  note: elementary's own document-open is a symlink to places/folder-open — this is the same art.
  also changes: fileopen (16, 22, 24), folder-open (16, 22, 24), folder-saved-search (16, 22, 24), gtk-directory (16, 22, 24), gtk-open (16, 22, 24), itmages-folder (16, 22, 24), project-open (16, 22, 24), sidebar-places (16, 22, 24), tag-folder (16, 22, 24)
  Comment: needs color change to blue according to theme
- [x] `document-properties` — elementary overlay
  names: document-properties, document-properties-symbolic · category: actions
  16 ← actions/16/document-properties.svg
  24, 22 ← actions/24/document-properties.svg
  48 ← actions/48/document-properties.svg
  document-properties-symbolic: symbolic ← actions/symbolic/document-properties-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16, 24, 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: bonobo-component-browser (16, 22, 24), games-config-custom (16, 22, 24), image-loading (16, 22, 24), items-tags (16, 22, 24), kstars_advanced (16, 22, 24), practice-setup (16, 22, 24), rabbitvcs-properties (16, 22, 24), settings-configure (16, 22, 24), stock_file-properties (16, 22, 24), tag-properties (16, 22, 24), xsi-document-properties-symbolic (symbolic)
- [x] `document-export` — elementary overlay
  names: document-export, document-export-symbolic, document-export-rtl, document-export-symbolic-rtl · category: actions
  16 ← actions/16/document-export.svg
  24, 22 ← actions/24/document-export.svg
  32 ← actions/32/document-export.svg
  48 ← actions/48/document-export.svg
  document-export-symbolic: symbolic ← actions/symbolic/document-export-symbolic.svg (recolor to ColorScheme-Text)
  document-export-rtl: 16 ← actions/16/document-export-rtl.svg
  document-export-rtl: 24, 22 ← actions/24/document-export-rtl.svg
  document-export-rtl: 32 ← actions/32/document-export-rtl.svg
  document-export-rtl: 48 ← actions/48/document-export-rtl.svg
  document-export-symbolic-rtl: symbolic ← actions/symbolic/document-export-symbolic-rtl.svg (recolor to ColorScheme-Text)
  also changes: document-export-ocal (16, 22, 24), document-export-table (16, 22, 24), document-import-symbolic-rtl (symbolic), export-symbolic (symbolic), rabbitvcs-export (16, 22, 24), stock_task-assigned-to (16, 22, 24)
- [x] `document-import` — elementary overlay
  names: document-import, document-import-symbolic, document-import-rtl, document-import-symbolic-rtl · category: actions
  16 ← actions/16/document-import.svg
  24, 22 ← actions/24/document-import.svg
  32 ← actions/32/document-import.svg
  48 ← actions/48/document-import.svg
  document-import-symbolic: symbolic ← actions/symbolic/document-import-symbolic.svg (recolor to ColorScheme-Text)
  document-import-rtl: 16 ← actions/16/document-import-rtl.svg
  document-import-rtl: 24, 22 ← actions/24/document-import-rtl.svg
  document-import-rtl: 32 ← actions/32/document-import-rtl.svg
  document-import-rtl: 48 ← actions/48/document-import-rtl.svg
  document-import-symbolic-rtl: symbolic ← actions/symbolic/document-import-symbolic-rtl.svg (recolor to ColorScheme-Text)
  also changes: albumfolder-importimages (16, 22, 24), document-export-symbolic-rtl (symbolic), document-import-ocal (16, 22, 24), rabbitvcs-import (16, 22, 24), stock_task-assigned (16, 22, 24)
- [x] `folder-new` — elementary overlay
  names: folder-new, folder-new-symbolic · category: actions
  16 ← actions/16/folder-new.svg
  24, 22 ← actions/24/folder-new.svg
  32 ← actions/32/folder-new.svg
  48 ← actions/48/folder-new.svg
  folder-new-symbolic: symbolic ← actions/symbolic/folder-new-symbolic.svg (recolor to ColorScheme-Text)
  also changes: action-albumfolder-importdir2 (16, 22, 24), add-folder-to-archive (16, 22, 24), albumfolder-importdir (16, 22, 24), albumfolder-new (16, 22, 24), archive-insert-directory (16, 22, 24), folder-add (16), folder_new (16, 22, 24), project_archive (16, 22, 24), stock_new-dir (16, 22, 24), xsi-folder-new-symbolic (symbolic)
  Comment: needs color change to blue according to theme
- [x] `tab-new` — elementary overlay from `window-new`
  names: tab-new · category: actions
  16 ← actions/16/window-new.svg
  24, 22 ← actions/24/window-new.svg
  32 ← actions/32/window-new.svg
  48 ← actions/48/window-new.svg
  note: New Tab (Dolphin, Konsole, Kate). Alternative: list-add.
  also changes: gtk-tab-new (22, 24), stock_new-tab (16, 22, 24), tab (16, 22, 24), tab_new (16, 22, 24)
- [x] `document-save-as` — elementary overlay
  names: document-save-as, document-save-as-symbolic · category: actions
  16 ← actions/16/document-save-as.svg
  24, 22 ← actions/24/document-save-as.svg
  32 ← actions/32/document-save-as.svg
  48 ← actions/48/document-save-as.svg
  document-save-as-symbolic: symbolic ← actions/symbolic/document-save-as-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: gtk-save-as (16, 22, 24), stock_save-as (16, 22, 24), xsi-document-save-as-symbolic (symbolic)
  Comment: needs color change to blue tones acording to theme
- [x] `bookmark-new` — elementary overlay
  names: bookmark-new, bookmark-new-symbolic · category: actions
  16 ← actions/16/bookmark-new.svg
  24, 22 ← actions/24/bookmark-new.svg
  32 ← actions/32/bookmark-new.svg
  48 ← actions/48/bookmark-new.svg
  bookmark-new-symbolic: symbolic ← actions/symbolic/bookmark-new-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: bookmark-add-folder (16, 22, 24), bookmark-add-symbolic (symbolic), bookmark-edit (16, 22, 24), bookmark-toolbar (16, 22, 24), bookmark_add (16, 22, 24), bookmarks-organize (16, 22, 24), bookmarks_list_add (16, 22, 24), poedit-status-bookmark (16, 22, 24), stock_add-bookmark (16, 22, 24), stock_help-add-bookmark (16, 22, 24), vmware-marker-add (16, 22, 24), xsi-bookmark-new-symbolic (symbolic)
- [x] `window-new` — elementary overlay
  names: window-new, window-new-symbolic · category: actions
  16 ← actions/16/window-new.svg
  24, 22 ← actions/24/window-new.svg
  32 ← actions/32/window-new.svg
  48 ← actions/48/window-new.svg
  window-new-symbolic: symbolic ← actions/symbolic/window-new-symbolic.svg (recolor to ColorScheme-Text)
  also changes: fm-new (16, 22, 24), window_new (16, 22, 24), xsi-window-new-symbolic (symbolic)
- [x] `document-revert` — elementary overlay
  names: document-revert, document-revert-symbolic, document-revert-rtl, document-revert-symbolic-rtl · category: actions
  16 ← actions/16/document-revert.svg
  24, 22 ← actions/24/document-revert.svg
  32 ← actions/32/document-revert.svg
  48 ← actions/48/document-revert.svg
  document-revert-symbolic: symbolic ← actions/symbolic/document-revert-symbolic.svg (recolor to ColorScheme-Text)
  document-revert-rtl: 16 ← actions/16/document-revert-rtl.svg
  document-revert-rtl: 24, 22 ← actions/24/document-revert-rtl.svg
  document-revert-rtl: 32 ← actions/32/document-revert-rtl.svg
  document-revert-rtl: 48 ← actions/48/document-revert-rtl.svg
  document-revert-symbolic-rtl: symbolic ← actions/symbolic/document-revert-symbolic-rtl.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: gtk-revert-to-saved-ltr (16, 22, 24), gtk-revert-to-saved-rtl (16, 22, 24), rabbitvcs-reset (16, 22, 24), revert (16, 22, 24), xsi-document-revert-symbolic (symbolic), xsi-document-revert-symbolic-rtl (symbolic)
- [x] `folder-open` — elementary overlay
  names: folder-open, folder-open-symbolic · category: actions
  16 ← places/16/folder-open.svg
  24, 22 ← places/24/folder-open.svg
  32 ← places/32/folder-open.svg
  48 ← places/48/folder-open.svg
  folder-open-symbolic: symbolic ← status/symbolic/folder-open-symbolic.svg (recolor to ColorScheme-Text)
  Comment: needs color change to blue tones according to theme
- [x] `trash-empty` — elementary overlay from `user-trash`
  names: trash-empty · category: actions
  16 ← places/16/user-trash.svg
  24, 22 ← places/24/user-trash.svg
  32 ← places/32/user-trash.svg
  48 ← places/48/user-trash.svg
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  note: KDE 'Move to Trash' and Dolphin 'Empty Trash'. Same art as the existing places/trash overlay's user-trash.
  note: elementary art at 16, 24, 32, 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
- [x] `document-print` — elementary overlay
  names: document-print, document-print-symbolic · category: actions
  24, 16, 22 ← actions/24/document-print.svg (→ devices/24/printer.svg)
  32 ← actions/32/document-print.svg (→ devices/32/printer.svg)
  48 ← actions/48/document-print.svg (→ devices/48/printer.svg)
  document-print-symbolic: symbolic ← actions/symbolic/document-print-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: document-print-direct (16, 22, 24), document-print-frame (16, 22, 24), dvips (16, 22, 24), fileprint (16, 22, 24), gtk-print (16, 22, 24), pdftex (16, 22, 24), stock_print (16, 22, 24), viewps (16, 22, 24), wizard_pstools (16, 22, 24), xsi-document-print-symbolic (symbolic)
- [x] `document-print-preview` — elementary overlay
  names: document-print-preview, document-print-preview-symbolic · category: actions
  24, 16, 22 ← actions/24/document-print-preview.svg
  32 ← actions/32/document-print-preview.svg
  48 ← actions/48/document-print-preview.svg
  document-print-preview-symbolic: symbolic ← actions/symbolic/document-print-preview-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: filequickprint (16, 22, 24), gtk-print-preview (16, 22, 24), stock_print-preview (16, 22, 24), xsi-document-print-preview-symbolic (symbolic)
- [x] `tag-new` — elementary overlay
  names: tag-new, tag-new-symbolic · category: actions
  16 ← actions/16/tag-new.svg
  24, 22 ← actions/24/tag-new.svg
  32 ← actions/32/tag-new.svg
  48 ← actions/48/tag-new.svg
  tag-new-symbolic: symbolic ← actions/symbolic/tag-new-symbolic.svg (recolor to ColorScheme-Text)
  also changes: gtg-tag-new (16, 22, 24)
- [x] `appointment-new` — elementary overlay
  names: appointment-new, appointment-new-symbolic · category: actions
  16 ← actions/16/appointment-new.svg
  24, 22 ← actions/24/appointment-new.svg
  32 ← actions/32/appointment-new.svg
  48 ← actions/48/appointment-new.svg
  appointment-new-symbolic: symbolic ← actions/symbolic/appointment-new-symbolic.svg (recolor to ColorScheme-Text)
  also changes: accept_time_event (16, 22, 24), amarok_clock (16, 22, 24), application-add-symbolic (symbolic), appointment (16, 22, 24), change-date-symbolic (symbolic), clock (16, 22, 24), clock-large (16, 22, 24), edit-undo-history (16, 22, 24), fm-new (16, 22, 24), gnumeric-pagesetup-hf-time (22, 24), knotes_alarm (16, 22, 24), kt-bandwidth-scheduler (16, 22, 24), shallow-history (16, 22, 24), smallclock (16, 22, 24), stock_new-appointment (16, 22, 24), window-new (16, 22, 24), window-new-symbolic (symbolic), window_new (16, 22, 24), xsi-appointment-new-symbolic (symbolic), xsi-window-new-symbolic (symbolic)
- [x] `address-book-new` — elementary overlay
  names: address-book-new, address-book-new-symbolic · category: actions
  16 ← actions/16/address-book-new.svg
  24, 22 ← actions/24/address-book-new.svg
  32 ← actions/32/address-book-new.svg
  48 ← actions/48/address-book-new.svg
  address-book-new-symbolic: symbolic ← actions/symbolic/address-book-new-symbolic.svg (recolor to ColorScheme-Text)
  also changes: addressbook-details (16, 22, 24), adress-book-new (16, 22, 24), bibtex (16, 22, 24), component (16, 22, 24), diag_component (16, 22, 24), documentation (16, 22, 24), file-catalog-symbolic (symbolic), identity (16, 22, 24), stock_new-address-book (16, 22, 24), tag-addressbook (16, 22, 24), umbrello_diagram_component (16, 22, 24), viewbib (16, 22, 24), xsi-address-book-new-symbolic (symbolic)
- [x] `contact-new` — elementary overlay
  names: contact-new, contact-new-symbolic · category: actions
  16 ← actions/16/contact-new.svg
  24, 22 ← actions/24/contact-new.svg
  32 ← actions/32/contact-new.svg
  48 ← actions/48/contact-new.svg
  contact-new-symbolic: symbolic ← actions/symbolic/contact-new-symbolic.svg (recolor to ColorScheme-Text)
  also changes: list-add-user (16, 22, 24), list-resource-add (16, 22, 24), stock_new-bcard (16, 22, 24), xsi-contact-new-symbolic (symbolic)
- [x] `document-page-setup` — elementary overlay
  names: document-page-setup, document-page-setup-symbolic · category: actions
  16 ← actions/16/document-page-setup.svg
  24, 22 ← actions/24/document-page-setup.svg
  32 ← actions/32/document-page-setup.svg
  48 ← actions/48/document-page-setup.svg
  document-page-setup-symbolic: symbolic ← actions/symbolic/document-page-setup-symbolic.svg (recolor to ColorScheme-Text)
  also changes: stock_print-setup (16, 22, 24), xsi-document-page-setup-symbolic (symbolic)
- [x] `folder-copy` — elementary overlay
  names: folder-copy · category: actions
  16 ← actions/16/folder-copy.svg
  24, 22 ← actions/24/folder-copy.svg
  also changes: stock_folder-copy (16, 22, 24)
  Comment: needs color change to blue tones according to theme
- [x] `folder-move` — elementary overlay
  names: folder-move · category: actions
  16 ← actions/16/folder-move.svg
  24, 22 ← actions/24/folder-move.svg
  also changes: stock_folder-move (16, 22, 24)
  Comment: needs color change to blue tones according to theme
- [x] `folder-tag` — elementary overlay
  names: folder-tag · category: actions
  16 ← places/16/folder-tag.svg
  24, 22 ← places/24/folder-tag.svg
  32 ← places/32/folder-tag.svg
  48 ← places/48/folder-tag.svg
- [x] `office-database-new` — elementary overlay
  names: office-database-new · category: actions
  16 ← actions/16/office-database-new.svg
  24, 22 ← actions/24/office-database-new.svg
  32 ← actions/32/office-database-new.svg
  48 ← actions/48/office-database-new.svg
- [x] `send-to` — elementary overlay
  names: send-to, send-to-symbolic · category: actions
  16 ← actions/16/send-to.svg
  24, 22 ← actions/24/send-to.svg
  send-to-symbolic: symbolic ← actions/symbolic/send-to-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16, 24px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: xsi-send-to-symbolic (symbolic)
- [x] `list-add` — elementary overlay
  names: list-add, list-add-symbolic · category: actions
  16 ← actions/16/list-add.svg
  24, 22 ← actions/24/list-add.svg
  32 ← actions/32/list-add.svg
  48 ← actions/48/list-add.svg
  list-add-symbolic: symbolic ← actions/symbolic/list-add-symbolic.svg (recolor to ColorScheme-Text)
  also changes: add (16, 22, 24), bqm-add (16, 22, 24), bqm-addqueue (16, 22, 24), edit-add (16, 22, 24), group (16, 22, 24), gtk-add (16, 22, 24), gtk-new (16, 22, 24), hb-add-queue (16, 22, 24), itmages-add (16, 22, 24), kdenlive-add-clip (16, 22, 24), kdenlive-add-color-clip (16, 22, 24), kdenlive-add-slide-clip (16, 22, 24), kdenlive-add-text-clip (16, 22, 24), kdenlive-insert-rect (16, 22, 24), kdenlive-insert-unicode (16, 22, 24), kt-add-feeds (16, 22, 24), lighttableadd (16, 22, 24), media-track-add-amarok (16, 22, 24), meeting-attending (16, 22, 24), menu_new (16, 22, 24), project_add (16, 22, 24), rabbitvcs-add (16, 22, 24), resource-group (16, 22, 24), stock_add (16), system-users (16, 22, 24), user-group-new (16, 22, 24), value-increase-symbolic (symbolic), xml-element-new (16, 22, 24), xsi-list-add-symbolic (symbolic), xsi-value-increase-symbolic (symbolic)
- [x] `edit-delete` — elementary overlay
  names: edit-delete, edit-delete-symbolic · category: actions
  16 ← actions/16/edit-delete.svg (→ places/16/user-trash.svg)
  24, 22 ← actions/24/edit-delete.svg (→ places/24/user-trash.svg)
  32 ← actions/32/edit-delete.svg (→ places/32/user-trash.svg)
  48 ← actions/48/edit-delete.svg (→ places/48/user-trash.svg)
  edit-delete-symbolic: symbolic ← actions/symbolic/edit-delete-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  note: elementary art at 16, 24, 32, 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: albumfolder-user-trash (16, 22, 24), bqm-remove (16, 22, 24), bqm-rmqueue (16, 22, 24), edit-delete-shred (16, 22, 24), editdelete (16, 22, 24), gtk-delete (16, 22, 24), hb-remove (16, 22, 24), itmages-remove (16, 22, 24), knotes_delete (16, 22, 24), kt-remove (16, 22, 24), rabbitvcs-delete (16, 22, 24), stock_delete (16, 22, 24), stock_not (16, 22, 24), tag-delete (16, 22, 24), usermenu-delete (16), xml-node-delete (16, 22, 24), xsi-edit-delete-symbolic (symbolic)
- [x] `edit-copy` — elementary overlay
  names: edit-copy, edit-copy-symbolic · category: actions
  16 ← actions/16/edit-copy.svg
  24, 22 ← actions/24/edit-copy.svg
  32 ← actions/32/edit-copy.svg
  48 ← actions/48/edit-copy.svg
  edit-copy-symbolic: symbolic ← actions/symbolic/edit-copy-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  note: elementary art at 16, 24, 32, 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: document-duplicate (16), document-multiple (16, 22, 24), editcopy (16, 22, 24), entry-clone (16, 22, 24), gtk-copy (16, 22, 24), gtk-dnd-multiple (16, 22, 24), msgid2msgstr (16, 22, 24), rabbitvcs-changes (16, 22, 24), specific-setup (16, 22, 24), stock_copy (16, 22, 24), stock_gtk-dnd-multiple (16, 22, 24), xsi-edit-copy-symbolic (symbolic)
- [x] `list-remove` — elementary overlay
  names: list-remove, list-remove-symbolic · category: actions
  16 ← actions/16/list-remove.svg
  24, 22 ← actions/24/list-remove.svg
  32 ← actions/32/list-remove.svg
  48 ← actions/48/list-remove.svg
  list-remove-symbolic: symbolic ← actions/symbolic/list-remove-symbolic.svg (recolor to ColorScheme-Text)
  also changes: gnumeric-object-line (22, 24), gtk-remove (16, 22, 24), kdenlive-deleffect (16, 22, 24), kt-remove-feeds (16, 22, 24), media-track-remove-amarok (16, 22, 24), menu_new_sep (16, 22, 24), no-marker (16, 22, 24), package-remove (16), project_remove (16, 22, 24), rabbitvcs-unstage (16, 22, 24), remove (16, 22, 24), remove-amarok (16, 22, 24), stock_insert-rule (16, 22, 24), stock_remove (16), value-decrease-symbolic (symbolic), xml-attribute-delete (16, 22, 24), xsi-list-remove-symbolic (symbolic), xsi-value-decrease-symbolic (symbolic)
- [x] `edit-undo` — elementary overlay
  names: edit-undo, edit-undo-symbolic, edit-undo-rtl, edit-undo-symbolic-rtl · category: actions
  16 ← actions/16/edit-undo.svg
  24, 22 ← actions/24/edit-undo.svg
  32 ← actions/32/edit-undo.svg
  48 ← actions/48/edit-undo.svg
  edit-undo-symbolic: symbolic ← actions/symbolic/edit-undo-symbolic.svg (recolor to ColorScheme-Text)
  edit-undo-rtl: 16 ← actions/16/edit-undo-rtl.svg
  edit-undo-rtl: 24, 22 ← actions/24/edit-undo-rtl.svg
  edit-undo-rtl: 32 ← actions/32/edit-undo-rtl.svg
  edit-undo-rtl: 48 ← actions/48/edit-undo-rtl.svg
  edit-undo-symbolic-rtl: symbolic ← actions/symbolic/edit-undo-symbolic-rtl.svg (→ actions/symbolic/edit-redo-symbolic.svg; recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: edit-redo-rtl-symbolic (symbolic), edit-redo-symbolic-rtl (symbolic), gtg-task-undismiss (16, 22, 24), gtk-undelete-ltr (16, 22, 24), gtk-undo-ltr (16), stock_undo (16, 22, 24), undo (16, 22, 24), xsi-edit-redo-symbolic-rtl (symbolic), xsi-edit-undo-symbolic (symbolic), xsi-edit-undo-symbolic-rtl (symbolic)
- [x] `edit-paste` — elementary overlay
  names: edit-paste, edit-paste-symbolic · category: actions
  16 ← actions/16/edit-paste.svg
  24, 22 ← actions/24/edit-paste.svg
  32 ← actions/32/edit-paste.svg
  48 ← actions/48/edit-paste.svg
  edit-paste-symbolic: symbolic ← actions/symbolic/edit-paste-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: clipman-symbolic (symbolic), editpaste (16, 22, 24), gtk-paste (16, 22, 24), rabbitvcs-editconflicts (16, 22, 24), stock_paste (16, 22, 24), xsi-edit-paste-symbolic (symbolic)
- [x] `edit-find` — elementary overlay
  names: edit-find, edit-find-symbolic · category: actions
  16 ← actions/16/edit-find.svg
  24, 22 ← actions/24/edit-find.svg
  48 ← actions/48/edit-find.svg
  edit-find-symbolic: symbolic ← actions/symbolic/edit-find-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: cm_search (16, 22, 24), document-edit-verify (16, 22, 24), document-preview-archive (16, 22, 24), edit-find-mail (16, 22, 24), edit-find-project (16, 22, 24), file-search-symbolic (symbolic), filefind (16, 22, 24), filegrep (16, 22, 24), find (16, 22, 24), gnome-searchtool (16, 22, 24), gtk-find (16, 22, 24), kfind (16, 22, 24), lens (16, 22, 24), preview (16, 22, 24), project_show (16, 22, 24), projectgrep (16, 22, 24), search (16, 22, 24), search-icon (16, 22, 24), search-symbolic (symbolic), spacefm-find (16, 22, 24), stock_search (16, 22, 24), system-search (16, 22, 24), system-search-symbolic (symbolic), transsearch (16, 22, 24), wizard_pdftools (16, 22, 24), xsi-edit-find-symbolic (symbolic), xsi-search-symbolic (symbolic), zoom (16, 22, 24)
- [x] `tools-check-spelling` — elementary overlay
  names: tools-check-spelling, tools-check-spelling-symbolic · category: actions
  16 ← actions/16/tools-check-spelling.svg
  24, 22 ← actions/24/tools-check-spelling.svg
  tools-check-spelling-symbolic: symbolic ← actions/symbolic/tools-check-spelling-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: check-spelling-symbolic (symbolic), gtk-spell-check (16, 22, 24), stock_spellcheck (16), xsi-check-spelling-symbolic (symbolic)
- [x] `edit-cut` — elementary overlay
  names: edit-cut, edit-cut-symbolic · category: actions
  16 ← actions/16/edit-cut.svg
  24, 22 ← actions/24/edit-cut.svg
  edit-cut-symbolic: symbolic ← actions/symbolic/edit-cut-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: editcut (16, 22, 24), gtk-cut (16, 22, 24), stock_cut (16, 22, 24), xsi-edit-cut-symbolic (symbolic)
- [x] `edit-find-replace` — elementary overlay
  names: edit-find-replace, edit-find-replace-symbolic · category: actions
  16 ← actions/16/edit-find-replace.svg
  24, 22 ← actions/24/edit-find-replace.svg
  48 ← actions/48/edit-find-replace.svg
  edit-find-replace-symbolic: symbolic ← actions/symbolic/edit-find-replace-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: gtk-find-and-replace (16, 22, 24), stock_search-and-replace (16), xsi-edit-find-replace-symbolic (symbolic)
- [x] `edit` — elementary overlay
  names: edit, edit-symbolic · category: actions
  16 ← actions/16/edit.svg (→ actions/16/document-edit.svg)
  24, 22 ← actions/24/edit.svg (→ actions/24/document-edit.svg)
  48 ← actions/48/edit.svg (→ actions/48/document-edit.svg)
  edit-symbolic: symbolic ← actions/symbolic/edit-symbolic.svg (→ actions/symbolic/document-edit-symbolic.svg; recolor to ColorScheme-Text)
  note: usage unknown: the name is a plain word, so the binary scan cannot tell icon lookups from other strings
  also changes: cell_edit (16, 22, 24), cm_multirename (16, 22, 24), document-edit (16, 22, 24), document-edit-symbolic (symbolic), document-sign (16, 22, 24), draw-freehand (16, 22, 24), edit-entry (16, 22, 24), edit-map (16, 22, 24), edit-rename (16, 22, 24), editimage (16, 22, 24), editor (16, 22, 24), edittext (16, 22, 24), entry-edit (16, 22, 24), gnumeric-comment-edit (16, 22, 24), gnumeric-link-edit (16, 22, 24), gtg-tags-all (16, 22, 24), gtk-edit (16, 22, 24), hb-edit (16, 22, 24), kdenlive-overwrite-edit (16, 22, 24), layer-rename (16, 22, 24), mail-task (16, 22, 24), media-track-edit-amarok (16, 22, 24), open-for-editing (16, 22, 24), rabbitvcs-rename (16, 22, 24), stock_edit (16, 22, 24), view-certificate-sign (16, 22, 24), view-time-schedule-edit (16, 22, 24), xsi-document-edit-symbolic (symbolic), xsi-edit-symbolic (symbolic)
- [x] `edit-flag` — elementary overlay
  names: edit-flag, edit-flag-symbolic · category: actions
  16 ← actions/16/edit-flag.svg
  24, 22 ← actions/24/edit-flag.svg
  edit-flag-symbolic: symbolic ← actions/symbolic/edit-flag-symbolic.svg (recolor to ColorScheme-Text)
- [x] `edit-mark` — elementary overlay
  names: edit-mark · category: actions
  24, 16, 22 ← actions/24/edit-mark.svg
- [x] `edit-undo-history` — elementary overlay
  names: edit-undo-history · category: actions
  16 ← actions/16/edit-undo-history.svg (→ places/16/folder-recent.svg)
  24, 22 ← actions/24/edit-undo-history.svg (→ places/24/folder-recent.svg)
- [x] `view-refresh` — elementary overlay
  names: view-refresh, view-refresh-symbolic · category: actions
  16 ← actions/16/view-refresh.svg
  24, 22 ← actions/24/view-refresh.svg
  32 ← actions/32/view-refresh.svg
  48 ← actions/48/view-refresh.svg
  view-refresh-symbolic: symbolic ← actions/symbolic/view-refresh-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: amarok_playlist_refresh (16, 22, 24), chronometer-reset (16, 22, 24), cm_refresh (16, 22, 24), collection-rescan-amarok (16, 22, 24), gtk-refresh (16, 22, 24), kt-change-tracker (16, 22, 24), poedit-update (16, 22, 24), rabbitvcs-refresh (16, 22, 24), refactor (16, 22, 24), refreshstructure (16, 22, 24), reload (16, 22, 24), reload3 (16, 22, 24), reload_all_tabs (16, 22, 24), reload_page (16, 22, 24), stock_refresh (16, 22, 24), stock_task-recurring (16, 22, 24), vm-restart (16, 22, 24), xsi-view-refresh-symbolic (symbolic)
- [x] `view-reader` — elementary overlay
  names: view-reader, view-reader-symbolic · category: actions
  16 ← actions/16/view-reader.svg
  24, 22 ← actions/24/view-reader.svg
  32 ← actions/32/view-reader.svg
  48 ← actions/48/view-reader.svg
  view-reader-symbolic: symbolic ← actions/symbolic/view-reader-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16, 24, 32, 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
- [x] `go-jump` — elementary overlay
  names: go-jump, go-jump-symbolic, go-jump-rtl, go-jump-symbolic-rtl · category: actions
  16 ← actions/16/go-jump.svg
  24, 22 ← actions/24/go-jump.svg
  32 ← actions/32/go-jump.svg
  48 ← actions/48/go-jump.svg
  go-jump-symbolic: symbolic ← actions/symbolic/go-jump-symbolic.svg (recolor to ColorScheme-Text)
  go-jump-rtl: 16 ← actions/16/go-jump-rtl.svg
  go-jump-rtl: 24, 22 ← actions/24/go-jump-rtl.svg
  go-jump-rtl: 32 ← actions/32/go-jump-rtl.svg
  go-jump-rtl: 48 ← actions/48/go-jump-rtl.svg
  go-jump-symbolic-rtl: symbolic ← actions/symbolic/go-jump-symbolic-rtl.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: go-jump-locationbar (16, 22, 24), gtk-jump-to-ltr (16, 22, 24), gtk-jump-to-rtl (16, 22, 24), kr_jumpback (16, 22, 24), rabbitvcs-switch (16, 22, 24), xsi-go-jump-symbolic (symbolic), xsi-go-jump-symbolic-rtl (symbolic)
- [x] `go-first-rtl` — elementary overlay
  names: go-first-rtl · category: actions
  16 ← actions/16/go-first-rtl.svg (→ actions/16/go-last.svg)
  24, 22 ← actions/24/go-first-rtl.svg (→ actions/24/go-last.svg)
  32 ← actions/32/go-first-rtl.svg (→ actions/32/go-last.svg)
  48 ← actions/48/go-first-rtl.svg (→ actions/48/go-last.svg)
- [x] `go-last-rtl` — elementary overlay
  names: go-last-rtl · category: actions
  16 ← actions/16/go-last-rtl.svg (→ actions/16/go-first.svg)
  24, 22 ← actions/24/go-last-rtl.svg (→ actions/24/go-first.svg)
  32 ← actions/32/go-last-rtl.svg (→ actions/32/go-first.svg)
  48 ← actions/48/go-last-rtl.svg (→ actions/48/go-first.svg)
- [x] `go-next-rtl` — elementary overlay
  names: go-next-rtl · category: actions
  16 ← actions/16/go-next-rtl.svg (→ actions/16/go-previous.svg)
  24, 22 ← actions/24/go-next-rtl.svg (→ actions/24/go-previous.svg)
  32 ← actions/32/go-next-rtl.svg (→ actions/32/go-previous.svg)
  48 ← actions/48/go-next-rtl.svg (→ actions/48/go-previous.svg)
- [x] `go-previous-rtl` — elementary overlay
  names: go-previous-rtl · category: actions
  16 ← actions/16/go-previous-rtl.svg (→ actions/16/go-next.svg)
  24, 22 ← actions/24/go-previous-rtl.svg (→ actions/24/go-next.svg)
  32 ← actions/32/go-previous-rtl.svg (→ actions/32/go-next.svg)
  48 ← actions/48/go-previous-rtl.svg (→ actions/48/go-next.svg)
- [x] `dialog-cancel` — elementary overlay from `process-stop`
  names: dialog-cancel · category: actions
  16 ← actions/16/process-stop.svg
  24, 22 ← actions/24/process-stop.svg
  32 ← actions/32/process-stop.svg
  48 ← actions/48/process-stop.svg
  note: Cancel buttons in KDE dialogs. Tela makes this an alias of window-close, so without this item it follows whatever window-close gets.
  also changes: cards-block (16, 22, 24), edit-none (16, 22, 24), end_of_life (16, 22, 24), filename-ignore-amarok (16, 22, 24)
- [x] `dialog-information` — elementary overlay
  names: dialog-information, dialog-information-symbolic · category: actions
  16 ← status/16/dialog-information.svg
  24, 22 ← status/24/dialog-information.svg
  32 ← status/32/dialog-information.svg
  48 ← status/48/dialog-information.svg
  dialog-information-symbolic: symbolic ← status/symbolic/dialog-information-symbolic.svg (recolor to ColorScheme-Text)
- [x] `dialog-warning` — elementary overlay
  names: dialog-warning, dialog-warning-symbolic · category: actions
  16 ← status/16/dialog-warning.svg
  24, 22 ← status/24/dialog-warning.svg
  32 ← status/32/dialog-warning.svg
  48 ← status/48/dialog-warning.svg
  dialog-warning-symbolic: symbolic ← status/symbolic/dialog-warning-symbolic.svg (recolor to ColorScheme-Text)
  also changes: itmages-warning (16, 22, 24), task-attempt (16, 22, 24), task-attention (16, 22, 24)
- [x] `dialog-error` — elementary overlay
  names: dialog-error, dialog-error-symbolic · category: actions
  16 ← status/16/dialog-error.svg
  24, 22 ← status/24/dialog-error.svg
  32 ← status/32/dialog-error.svg
  48 ← status/48/dialog-error.svg
  dialog-error-symbolic: symbolic ← status/symbolic/dialog-error-symbolic.svg (recolor to ColorScheme-Text)
  also changes: SuggestionError (16, 22, 24), itmages-error (16, 22, 24), messagebox_warning (16, 22, 24)
- [x] `dialog-ok-apply` — elementary overlay from `process-completed`
  names: dialog-ok-apply · category: actions
  16 ← status/16/process-completed.svg
  24, 22 ← status/24/process-completed.svg
  32 ← status/32/process-completed.svg
  48 ← status/48/process-completed.svg
  note: OK/Apply buttons in KDE dialogs; dialog-ok is a Tela alias of this file, so it follows.
  also changes: Finished (16, 22, 24), answer (16, 22, 24), answer-correct (16, 22, 24), blueman-trust (16, 22, 24), check-filled (16, 22, 24), checkbox (16, 22, 24), checkmark (16, 22, 24), dialog-apply (16, 22, 24), dialog-ok (16, 22, 24), dialog-yes (16, 22, 24), enterprise (16, 22, 24), games-endturn (16, 22, 24), games-solve (16, 22, 24), gtk-apply (16, 22, 24), gtk-ok (16, 22, 24), gtk-yes (16, 22, 24), hb-complete (16, 22, 24), kr_mountman (16, 22, 24), kt-check-data (16, 22, 24), mail-mark-task (16, 22, 24), poedit-validate (16, 22, 24), preflight-verifier (16, 22, 24), stock_calc-accept (16, 22, 24), stock_check-filled (16, 22, 24), stock_mark (16, 22, 24), stock_todo (16, 22, 24), stock_yes (16, 22, 24)
- [x] `process-stop` — elementary overlay
  names: process-stop, process-stop-symbolic · category: actions
  16 ← actions/16/process-stop.svg
  24, 22 ← actions/24/process-stop.svg
  32 ← actions/32/process-stop.svg
  48 ← actions/48/process-stop.svg
  process-stop-symbolic: symbolic ← actions/symbolic/process-stop-symbolic.svg (recolor to ColorScheme-Text)
  also changes: blueman-untrust (16, 22, 24), gtk-cancel (16, 22, 24), gtk-stop (16, 22, 24), hb-stop (16, 22, 24), itmages-cancel (16, 22, 24), rabbitvcs-stop (16, 22, 24), stock_stop (16, 22, 24), stop (16, 22, 24), vm-tools-install-cancel (16, 22, 24), xsi-process-stop-symbolic (symbolic)
- [x] `help-contents` — elementary overlay
  names: help-contents, help-contents-symbolic · category: actions
  16 ← actions/16/help-contents.svg (→ categories/16/system-help.svg)
  24, 22 ← actions/24/help-contents.svg (→ categories/24/system-help.svg)
  32 ← actions/32/help-contents.svg (→ categories/32/system-help.svg)
  48 ← actions/48/help-contents.svg (→ categories/48/system-help.svg)
  help-contents-symbolic: symbolic ← actions/symbolic/help-contents-symbolic.svg (→ categories/symbolic/system-help-symbolic.svg; recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: contexthelp (16, 22, 24), help (16, 22, 24), help-contents-vmware (16, 22, 24), help-contextual (16, 22, 24), itmages-question (16, 22, 24), lpi-help (16, 22, 24), rabbitvcs-help (16, 22, 24), system-help (16, 22, 24), xsi-help-contents-symbolic (symbolic)
- [x] `dialog-password` — elementary overlay
  names: dialog-password, dialog-password-symbolic · category: actions
  24, 16, 22 ← status/24/dialog-password.svg
  32 ← status/32/dialog-password.svg
  48 ← status/48/dialog-password.svg
  dialog-password-symbolic: symbolic ← status/symbolic/dialog-password-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 24, 32, 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
- [x] `system-lock-screen` — elementary overlay
  names: system-lock-screen, system-lock-screen-symbolic · category: actions
  32, 16, 22, 24 ← actions/32/system-lock-screen.svg
  48 ← actions/48/system-lock-screen.svg
  system-lock-screen-symbolic: symbolic ← actions/symbolic/system-lock-screen-symbolic.svg (recolor to ColorScheme-Text)
  also changes: xfsm-lock (16, 22, 24)
- [x] `system-log-out` — elementary overlay
  names: system-log-out, system-log-out-symbolic · category: actions
  32, 16, 22, 24 ← actions/32/system-log-out.svg
  48 ← actions/48/system-log-out.svg
  system-log-out-symbolic: symbolic ← actions/symbolic/system-log-out-symbolic.svg (recolor to ColorScheme-Text)
  also changes: gnome-logout (16), math1 (16), xfsm-logout (16, 22, 24)
- [x] `system-reboot` — elementary overlay
  names: system-reboot, system-reboot-symbolic · category: actions
  32, 16, 22, 24 ← actions/32/system-reboot.svg
  48 ← actions/48/system-reboot.svg
  system-reboot-symbolic: symbolic ← actions/symbolic/system-reboot-symbolic.svg (→ status/symbolic/rotation-allowed-symbolic.svg; recolor to ColorScheme-Text)
- [x] `system-search` — elementary overlay
  names: system-search, system-search-symbolic · category: actions
  16 ← actions/16/system-search.svg (→ actions/16/edit-find.svg)
  24, 22 ← actions/24/system-search.svg (→ actions/24/edit-find.svg)
  48 ← actions/48/system-search.svg (→ actions/48/edit-find.svg)
  system-search-symbolic: symbolic ← actions/symbolic/system-search-symbolic.svg (→ actions/symbolic/edit-find-symbolic.svg; recolor to ColorScheme-Text)
  also changes: kfind (16)
- [x] `system-shutdown` — elementary overlay
  names: system-shutdown, system-shutdown-symbolic · category: actions
  32, 16, 22, 24 ← actions/32/system-shutdown.svg
  48 ← actions/48/system-shutdown.svg
  system-shutdown-symbolic: symbolic ← actions/symbolic/system-shutdown-symbolic.svg (recolor to ColorScheme-Text)
  also changes: xfsm-shutdown (16, 22, 24)
- [x] `system-suspend` — elementary overlay
  names: system-suspend · category: actions
  32, 16, 22, 24 ← actions/32/system-suspend.svg
  also changes: system-suspend-symbolic (22, 32), xfsm-suspend (16, 22, 24)
- [x] `system-help` — elementary overlay
  names: system-help, system-help-symbolic · category: actions
  16 ← categories/16/system-help.svg
  24, 22 ← categories/24/system-help.svg
  32 ← categories/32/system-help.svg
  48 ← categories/48/system-help.svg
  system-help-symbolic: symbolic ← categories/symbolic/system-help-symbolic.svg (recolor to ColorScheme-Text)
- [x] `starred` — elementary overlay
  names: starred, starred-symbolic · category: actions
  16, 22, 24 ← status/16/starred.svg
  starred-symbolic: symbolic ← status/symbolic/starred-symbolic.svg (recolor to ColorScheme-Text)
  also changes: xsi-starred-symbolic (symbolic)
- [x] `applications-other` — elementary overlay
  names: applications-other, applications-other-symbolic · category: actions
  16 ← categories/16/applications-other.svg
  24, 22 ← categories/24/applications-other.svg
  32 ← categories/32/applications-other.svg
  48 ← categories/48/applications-other.svg
  applications-other-symbolic: symbolic ← categories/symbolic/applications-other-symbolic.svg (recolor to ColorScheme-Text)
- [x] `gnome-lockscreen` — elementary overlay
  names: gnome-lockscreen · category: actions
  48, 16, 22, 24 ← actions/48/gnome-lockscreen.svg (→ actions/48/system-lock-screen.svg)
- [x] `jockey-proprietary` — elementary overlay
  names: jockey-proprietary · category: actions
  24, 16, 22 ← actions/24/jockey-proprietary.svg (→ status/24/locked.svg)
- [x] `non-starred` — elementary overlay
  names: non-starred, non-starred-symbolic · category: actions
  16, 22, 24 ← status/16/non-starred.svg
  non-starred-symbolic: symbolic ← status/symbolic/non-starred-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: xsi-non-starred-symbolic (symbolic)
- [x] `off` — elementary overlay
  names: off, off-symbolic · category: actions
  16 ← status/16/off.svg (→ actions/16/process-stop.svg)
  24, 22 ← status/24/off.svg (→ actions/24/process-stop.svg)
  32 ← status/32/off.svg (→ actions/32/process-stop.svg)
  48 ← status/48/off.svg (→ actions/48/process-stop.svg)
  off-symbolic: symbolic ← status/symbolic/off-symbolic.svg (→ actions/symbolic/process-stop-symbolic.svg; recolor to ColorScheme-Text)
  note: usage unknown: the name is a plain word, so the binary scan cannot tell icon lookups from other strings
- [x] `off-outline` — elementary overlay
  names: off-outline, off-outline-symbolic · category: actions
  16 ← status/16/off-outline.svg (→ actions/16/process-stop.svg)
  24, 22 ← status/24/off-outline.svg (→ actions/24/process-stop.svg)
  32 ← status/32/off-outline.svg (→ actions/32/process-stop.svg)
  48 ← status/48/off-outline.svg (→ actions/48/process-stop.svg)
  off-outline-symbolic: symbolic ← status/symbolic/off-outline-symbolic.svg (recolor to ColorScheme-Text)
- [x] `on` — elementary overlay
  names: on, on-symbolic · category: actions
  16 ← status/16/on.svg (→ status/16/process-completed.svg)
  24, 22 ← status/24/on.svg (→ status/24/process-completed.svg)
  32 ← status/32/on.svg (→ status/32/process-completed.svg)
  48 ← status/48/on.svg (→ status/48/process-completed.svg)
  on-symbolic: symbolic ← status/symbolic/on-symbolic.svg (recolor to ColorScheme-Text)
  note: usage unknown: the name is a plain word, so the binary scan cannot tell icon lookups from other strings
- [x] `on-outline` — elementary overlay
  names: on-outline, on-outline-symbolic · category: actions
  16 ← status/16/on-outline.svg (→ status/16/process-completed.svg)
  24, 22 ← status/24/on-outline.svg (→ status/24/process-completed.svg)
  32 ← status/32/on-outline.svg (→ status/32/process-completed.svg)
  48 ← status/48/on-outline.svg (→ status/48/process-completed.svg)
  on-outline-symbolic: symbolic ← status/symbolic/on-outline-symbolic.svg (recolor to ColorScheme-Text)
- [x] `semi-starred` — elementary overlay
  names: semi-starred, semi-starred-symbolic, semi-starred-rtl, semi-starred-symbolic-rtl · category: actions
  16, 22, 24 ← status/16/semi-starred.svg
  semi-starred-symbolic: symbolic ← status/symbolic/semi-starred-symbolic.svg (recolor to ColorScheme-Text)
  semi-starred-rtl: 16, 22, 24 ← status/16/semi-starred-rtl.svg
  semi-starred-symbolic-rtl: symbolic ← status/symbolic/semi-starred-symbolic-rtl.svg (recolor to ColorScheme-Text)
  also changes: xsi-semi-starred-symbolic (symbolic), xsi-semi-starred-symbolic-rtl (symbolic)
- [x] `system-restart` — elementary overlay
  names: system-restart · category: actions
  24, 16, 22 ← actions/24/system-restart.svg (→ actions/24/view-refresh.svg)
  32 ← actions/32/system-restart.svg (→ actions/32/system-reboot.svg)
  48 ← actions/48/system-restart.svg (→ actions/48/system-reboot.svg)
  also changes: system-reboot (16, 22, 24, 32), xfsm-reboot (16, 22, 24)
- [x] `x-office-address-book` — elementary overlay
  names: x-office-address-book, x-office-address-book-symbolic · category: actions
  16 ← mimes/16/x-office-address-book.svg
  24, 22 ← mimes/24/x-office-address-book.svg
  32 ← mimes/32/x-office-address-book.svg
  48 ← mimes/48/x-office-address-book.svg
  x-office-address-book-symbolic: symbolic ← mimes/symbolic/x-office-address-book-symbolic.svg (recolor to ColorScheme-Text)
- [x] `call-start` — elementary overlay
  names: call-start, call-start-symbolic · category: actions
  16, 22, 24 ← actions/16/call-start.svg
  call-start-symbolic: symbolic ← actions/symbolic/call-start-symbolic.svg (recolor to ColorScheme-Text)
  also changes: xsi-call-start-symbolic (symbolic)
- [x] `call-end` — elementary overlay
  names: call-end, call-end-symbolic · category: actions
  16, 22, 24 ← actions/16/call-end.svg (→ actions/16/call-stop.svg)
  call-end-symbolic: symbolic ← actions/symbolic/call-end-symbolic.svg (→ actions/symbolic/call-stop-symbolic.svg; recolor to ColorScheme-Text)
- [x] `call-stop` — elementary overlay
  names: call-stop, call-stop-symbolic · category: actions
  16, 22, 24 ← actions/16/call-stop.svg
  call-stop-symbolic: symbolic ← actions/symbolic/call-stop-symbolic.svg (recolor to ColorScheme-Text)
  also changes: call-end-symbolic (symbolic), xsi-call-stop-symbolic (symbolic)
- [x] `mail-message-new` — elementary overlay
  names: mail-message-new, mail-message-new-symbolic · category: actions
  16 ← actions/16/mail-message-new.svg
  24, 22 ← actions/24/mail-message-new.svg
  mail-message-new-symbolic: symbolic ← actions/symbolic/mail-message-new-symbolic.svg (recolor to ColorScheme-Text)
  also changes: gnome-stock-mail-new (16, 22, 24), im-message-new (16, 22, 24), mail_new (16, 22, 24), stock_mail-compose (16, 22, 24), xsi-mail-message-new-symbolic (symbolic)
- [x] `mail-forward` — elementary overlay
  names: mail-forward, mail-forward-symbolic, mail-forward-rtl, mail-forward-symbolic-rtl · category: actions
  16 ← actions/16/mail-forward.svg
  24, 22 ← actions/24/mail-forward.svg
  32 ← actions/32/mail-forward.svg
  48 ← actions/48/mail-forward.svg
  mail-forward-symbolic: symbolic ← actions/symbolic/mail-forward-symbolic.svg (recolor to ColorScheme-Text)
  mail-forward-rtl: 16 ← actions/16/mail-forward-rtl.svg
  mail-forward-rtl: 24, 22 ← actions/24/mail-forward-rtl.svg
  mail-forward-rtl: 32 ← actions/32/mail-forward-rtl.svg
  mail-forward-rtl: 48 ← actions/48/mail-forward-rtl.svg
  mail-forward-symbolic-rtl: symbolic ← actions/symbolic/mail-forward-symbolic-rtl.svg (→ actions/symbolic/mail-reply-sender-symbolic.svg; recolor to ColorScheme-Text)
  also changes: gnome-stock-mail-fwd (16, 22, 24), mail-forwarded (16, 22, 24), mail-forwarded-replied (16, 22, 24), mail-reply-rtl-symbolic (symbolic), mail-reply-sender-rtl-symbolic (symbolic), mail-reply-sender-symbolic-rtl (symbolic), mail_forward (16, 22, 24), stock_mail-forward (16, 22, 24), umbr-coll-message-asynchronous (16, 22, 24), umbr-coll-message-synchronous (16, 22, 24), umbr-message-asynchronous (16, 22, 24), umbr-message-synchronous (16, 22, 24), xsi-mail-forward-symbolic (symbolic), xsi-mail-forward-symbolic-rtl (symbolic), xsi-mail-reply-sender-symbolic-rtl (symbolic)
- [x] `mail-forwarded` — elementary overlay
  names: mail-forwarded, mail-forwarded-symbolic, mail-forwarded-rtl, mail-forwarded-symbolic-rtl · category: actions
  16 ← status/16/mail-forwarded.svg (→ actions/16/mail-forward.svg)
  24, 22 ← status/24/mail-forwarded.svg (→ actions/24/mail-forward.svg)
  mail-forwarded-symbolic: symbolic ← status/symbolic/mail-forwarded-symbolic.svg (→ actions/symbolic/go-next-symbolic.svg; recolor to ColorScheme-Text)
  mail-forwarded-rtl: 16 ← status/16/mail-forwarded-rtl.svg (→ actions/16/mail-forward-rtl.svg)
  mail-forwarded-rtl: 24, 22 ← status/24/mail-forwarded-rtl.svg (→ actions/24/mail-forward-rtl.svg)
  mail-forwarded-symbolic-rtl: symbolic ← status/symbolic/mail-forwarded-symbolic-rtl.svg (→ actions/symbolic/go-previous-symbolic.svg; recolor to ColorScheme-Text)
  also changes: umbr-coll-message-synchronous (16, 22, 24), umbr-message-synchronous (16, 22, 24)
- [x] `mail-inbox` — elementary overlay
  names: mail-inbox, mail-inbox-symbolic · category: actions
  16, 22, 24 ← places/16/mail-inbox.svg
  mail-inbox-symbolic: symbolic ← places/symbolic/mail-inbox-symbolic.svg (recolor to ColorScheme-Text)
- [x] `mail-mark-important` — elementary overlay
  names: mail-mark-important, mail-mark-important-symbolic · category: actions
  16, 22, 24 ← actions/16/mail-mark-important.svg
  mail-mark-important-symbolic: symbolic ← actions/symbolic/mail-mark-important-symbolic.svg (recolor to ColorScheme-Text)
  also changes: xsi-mail-mark-important-symbolic (symbolic)
- [x] `mail-mark-junk` — elementary overlay
  names: mail-mark-junk, mail-mark-junk-symbolic · category: actions
  16 ← actions/16/mail-mark-junk.svg
  24, 22 ← actions/24/mail-mark-junk.svg
  32 ← actions/32/mail-mark-junk.svg
  48 ← actions/48/mail-mark-junk.svg
  mail-mark-junk-symbolic: symbolic ← actions/symbolic/mail-mark-junk-symbolic.svg (recolor to ColorScheme-Text)
  also changes: lpi-bug (16, 22, 24), mail-spam-symbolic (symbolic), mail_spam (16), rabbitvcs-bug (16, 22, 24), spam-symbolic (symbolic), stock_spam (16), umbr-message-lost (16, 22, 24), xsi-mail-mark-junk-symbolic (symbolic)
- [x] `mail-mark-notjunk` — elementary overlay
  names: mail-mark-notjunk, mail-mark-notjunk-symbolic · category: actions
  16 ← actions/16/mail-mark-notjunk.svg
  24, 22 ← actions/24/mail-mark-notjunk.svg
  32 ← actions/32/mail-mark-notjunk.svg
  48 ← actions/48/mail-mark-notjunk.svg
  mail-mark-notjunk-symbolic: symbolic ← actions/symbolic/mail-mark-notjunk-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16, 24, 32, 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: umbr-message-found (16, 22, 24), xsi-mail-mark-notjunk-symbolic (symbolic)
- [x] `mail-outbox` — elementary overlay
  names: mail-outbox, mail-outbox-symbolic · category: actions
  16, 22, 24 ← places/16/mail-outbox.svg
  mail-outbox-symbolic: symbolic ← places/symbolic/mail-outbox-symbolic.svg (recolor to ColorScheme-Text)
- [x] `mail-replied` — elementary overlay
  names: mail-replied, mail-replied-symbolic, mail-replied-rtl, mail-replied-symbolic-rtl · category: actions
  16 ← status/16/mail-replied.svg (→ actions/16/mail-reply-sender.svg)
  24, 22 ← status/24/mail-replied.svg (→ actions/24/mail-reply-sender.svg)
  mail-replied-symbolic: symbolic ← status/symbolic/mail-replied-symbolic.svg (→ actions/symbolic/mail-reply-sender-symbolic.svg; recolor to ColorScheme-Text)
  mail-replied-rtl: 16 ← status/16/mail-replied-rtl.svg (→ actions/16/mail-reply-sender-rtl.svg)
  mail-replied-rtl: 24, 22 ← status/24/mail-replied-rtl.svg (→ actions/24/mail-reply-sender-rtl.svg)
  mail-replied-symbolic-rtl: symbolic ← status/symbolic/mail-replied-symbolic-rtl.svg (→ actions/symbolic/mail-forward-symbolic.svg; recolor to ColorScheme-Text)
- [x] `mail-reply-all` — elementary overlay
  names: mail-reply-all, mail-reply-all-symbolic, mail-reply-all-rtl, mail-reply-all-symbolic-rtl · category: actions
  16 ← actions/16/mail-reply-all.svg
  24, 22 ← actions/24/mail-reply-all.svg
  32 ← actions/32/mail-reply-all.svg
  48 ← actions/48/mail-reply-all.svg
  mail-reply-all-symbolic: symbolic ← actions/symbolic/mail-reply-all-symbolic.svg (recolor to ColorScheme-Text)
  mail-reply-all-rtl: 16 ← actions/16/mail-reply-all-rtl.svg
  mail-reply-all-rtl: 24, 22 ← actions/24/mail-reply-all-rtl.svg
  mail-reply-all-rtl: 32 ← actions/32/mail-reply-all-rtl.svg
  mail-reply-all-rtl: 48 ← actions/48/mail-reply-all-rtl.svg
  mail-reply-all-symbolic-rtl: symbolic ← actions/symbolic/mail-reply-all-symbolic-rtl.svg (recolor to ColorScheme-Text)
  also changes: mail-reply-custom-all (16, 22, 24), mail-reply-list (16, 22, 24), mail_replyall (16, 22, 24), stock_mail-reply-to-all (16, 22, 24), xsi-mail-reply-all-symbolic (symbolic), xsi-mail-reply-all-symbolic-rtl (symbolic)
- [x] `mail-reply-sender` — elementary overlay
  names: mail-reply-sender, mail-reply-sender-symbolic, mail-reply-sender-rtl, mail-reply-sender-symbolic-rtl · category: actions
  16 ← actions/16/mail-reply-sender.svg
  24, 22 ← actions/24/mail-reply-sender.svg
  32 ← actions/32/mail-reply-sender.svg
  48 ← actions/48/mail-reply-sender.svg
  mail-reply-sender-symbolic: symbolic ← actions/symbolic/mail-reply-sender-symbolic.svg (recolor to ColorScheme-Text)
  mail-reply-sender-rtl: 16 ← actions/16/mail-reply-sender-rtl.svg
  mail-reply-sender-rtl: 24, 22 ← actions/24/mail-reply-sender-rtl.svg
  mail-reply-sender-rtl: 32 ← actions/32/mail-reply-sender-rtl.svg
  mail-reply-sender-rtl: 48 ← actions/48/mail-reply-sender-rtl.svg
  mail-reply-sender-symbolic-rtl: symbolic ← actions/symbolic/mail-reply-sender-symbolic-rtl.svg (→ actions/symbolic/mail-forward-symbolic.svg; recolor to ColorScheme-Text)
  also changes: gnome-stock-mail-rpl (16, 22, 24), kt-restore-defaults (16, 22, 24), mail-replied (16, 22, 24), mail-reply-custom (16, 22, 24), mail_reply (16, 22, 24), stock_mail-reply (16, 22, 24), xsi-mail-reply-sender-symbolic (symbolic), xsi-mail-reply-sender-symbolic-rtl (symbolic)
- [x] `mail-send-receive` — elementary overlay
  names: mail-send-receive, mail-send-receive-symbolic · category: actions
  16 ← actions/16/mail-send-receive.svg
  24, 22 ← actions/24/mail-send-receive.svg
  mail-send-receive-symbolic: symbolic ← actions/symbolic/mail-send-receive-symbolic.svg (recolor to ColorScheme-Text)
  also changes: stock_mail-send-receive (16, 22, 24), xsi-mail-send-receive-symbolic (symbolic)
- [x] `mail-unread` — elementary overlay
  names: mail-unread, mail-unread-symbolic · category: actions
  16, 22, 24 ← status/16/mail-unread.svg
  mail-unread-symbolic: symbolic ← status/symbolic/mail-unread-symbolic.svg (recolor to ColorScheme-Text)
- [x] `format-text-bold` — elementary overlay
  names: format-text-bold, format-text-bold-symbolic · category: actions
  16 ← actions/16/format-text-bold.svg
  24, 22 ← actions/24/format-text-bold.svg
  format-text-bold-symbolic: symbolic ← actions/symbolic/format-text-bold-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16, 24px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: gtk-bold (16, 22, 24), stock_text_bold (16), strong (16, 22, 24), text_bold (16), xsi-format-text-bold-symbolic (symbolic)
- [x] `format-text-strikethrough` — elementary overlay
  names: format-text-strikethrough, format-text-strikethrough-symbolic · category: actions
  16 ← actions/16/format-text-strikethrough.svg
  24, 22 ← actions/24/format-text-strikethrough.svg
  format-text-strikethrough-symbolic: symbolic ← actions/symbolic/format-text-strikethrough-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16, 24px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: gtk-strikethrough (16), stock_text-strikethrough (16), text_strike (16), xsi-format-text-strikethrough-symbolic (symbolic)
- [x] `format-text-underline` — elementary overlay
  names: format-text-underline, format-text-underline-symbolic · category: actions
  16 ← actions/16/format-text-underline.svg
  24, 22 ← actions/24/format-text-underline.svg
  format-text-underline-symbolic: symbolic ← actions/symbolic/format-text-underline-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16, 24px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: distribute-vertical-baseline (16, 22, 24), gtk-underline (16), stock_text_underlined (16), text_under (16), xsi-format-text-underline-symbolic (symbolic)
- [x] `format-text-italic` — elementary overlay
  names: format-text-italic, format-text-italic-symbolic · category: actions
  16 ← actions/16/format-text-italic.svg
  24, 22 ← actions/24/format-text-italic.svg
  format-text-italic-symbolic: symbolic ← actions/symbolic/format-text-italic-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16, 24px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: artistictext-tool (16, 22, 24), emph (16, 22, 24), gtk-italic (16, 22, 24), slanted (16, 22, 24), stock_text_italic (16), text_italic (16), xsi-format-text-italic-symbolic (symbolic)
- [x] `format-text-highlight` — elementary overlay
  names: format-text-highlight · category: actions
  24, 16, 22 ← actions/24/format-text-highlight.svg
- [x] `format-text-none` — elementary overlay
  names: format-text-none · category: actions
  16 ← actions/16/format-text-none.svg
  24, 22 ← actions/24/format-text-none.svg
- [x] `object-locked` — elementary overlay
  names: object-locked · category: actions
  16 ← actions/16/object-locked.svg (→ status/16/locked.svg)
  24, 22 ← actions/24/object-locked.svg (→ status/24/locked.svg)
  32 ← actions/32/object-locked.svg (→ status/32/locked.svg)
  48 ← actions/48/object-locked.svg (→ status/48/locked.svg)
  also changes: database-lock (16, 22, 24), document-encrypt (16, 22, 24), document-encrypted (16, 22, 24), gnome-lockscreen (16, 22, 24), gtk-authentication (16, 22, 24), kdenlive-lock (16, 22, 24), kdenlive-unlock (16, 22, 24), kt-encrypted (16, 22, 24), lock (16, 22, 24), preferences-desktop-user-password (16, 22, 24), rabbitvcs-lock (16, 22, 24), system-lock-screen (24), udiskie-lock (16, 22, 24), vmware-locked (16, 22, 24), xfsm-lock (24)
- [x] `object-unlocked` — elementary overlay
  names: object-unlocked · category: actions
  16 ← actions/16/object-unlocked.svg (→ status/16/changes-allow.svg)
  24, 22 ← actions/24/object-unlocked.svg (→ status/24/changes-allow.svg)
  32 ← actions/32/object-unlocked.svg (→ status/32/changes-allow.svg)
  48 ← actions/48/object-unlocked.svg (→ status/48/changes-allow.svg)
  also changes: document-decrypt (16, 22, 24), gpg (16, 22, 24), gpgsm (16, 22, 24), rabbitvcs-unlock (16, 22, 24), udiskie-unlock (16, 22, 24), unlock (16, 22, 24), vmware-unlocked (16, 22, 24)

## Rejected

- [x] `document-open` — elementary overlay
  names: document-open, document-open-symbolic · category: actions
  16 ← actions/16/document-open.svg (→ places/16/folder-open.svg)
  24, 22 ← actions/24/document-open.svg
  32 ← actions/32/document-open.svg (→ places/32/folder-open.svg)
  48 ← actions/48/document-open.svg (→ places/48/folder-open.svg)
  document-open-symbolic: symbolic ← actions/symbolic/document-open-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: document-open-data (16, 22, 24), hb-source (16, 22, 24), irc-join-channel (16, 22, 24), scriptopen (16, 22, 24), xsi-document-open-symbolic (symbolic)
- [x] `edit-clear` — elementary overlay
  names: edit-clear, edit-clear-symbolic · category: actions
  24, 16, 22 ← actions/24/edit-clear.svg
  32 ← actions/32/edit-clear.svg
  48 ← actions/48/edit-clear.svg
  edit-clear-symbolic: symbolic ← actions/symbolic/edit-clear-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: clear_left (16, 22, 24), document-cleanup (16, 22, 24), edit-clear-history (16, 22, 24), edit-clear-locationbar-rtl (16, 22, 24), editclear (16, 22, 24), xsi-edit-clear-symbolic (symbolic)
  Comment: actually i dont like this icon
- [ ] `edit-select-all` — elementary overlay
  names: edit-select-all, edit-select-all-symbolic · category: actions
  16 ← actions/16/edit-select-all.svg
  24, 22 ← actions/24/edit-select-all.svg
  32 ← actions/32/edit-select-all.svg
  48 ← actions/48/edit-select-all.svg
  edit-select-all-symbolic: symbolic ← actions/symbolic/edit-select-all-symbolic.svg (recolor to ColorScheme-Text)
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  also changes: gtk-select-all (16, 22, 24), itmages-select-all (16, 22, 24), region (16, 22, 24), stock_select-all (16, 22, 24), xsi-edit-select-all-symbolic (symbolic)
  Comment: needs better icon
- [x] `edit-clear-list` — elementary overlay from `edit-clear`
  names: edit-clear-list · category: actions
  24, 16, 22 ← actions/24/edit-clear.svg
  32 ← actions/32/edit-clear.svg
  48 ← actions/48/edit-clear.svg
  note: Clear list. edit-clear-all is a Tela alias of it.
  also changes: edit-clear-all (16, 22, 24), gtg-tags-none (16, 22, 24), gtk-clear (16, 22, 24), gtk-quit (16, 22, 24), itmages-quit (16, 22, 24), kdenlive-align-none (16, 22, 24)
  Comment: needs better icon
- [x] `edit-clear-rtl` — elementary overlay from `edit-clear`
  names: edit-clear-rtl · category: actions
  24, 16, 22 ← actions/24/edit-clear.svg
  32 ← actions/32/edit-clear.svg
  48 ← actions/48/edit-clear.svg
  note: Tela's mirrored edit-clear. edit-clear-locationbar-ltr (the location-bar clear button in LTR locales) is a Tela alias of it.
  also changes: edit-clear-locationbar-ltr (16, 22, 24), locationbar_erase (16, 22, 24)
  Comment: needs better
- [ ] `find-location` — elementary overlay
  names: find-location, find-location-symbolic · category: actions
  32, 16, 22, 24 ← actions/32/find-location.svg
  48 ← actions/48/find-location.svg
  find-location-symbolic: symbolic ← actions/symbolic/find-location-symbolic.svg (recolor to ColorScheme-Text)
  also changes: xsi-find-location-symbolic (symbolic)
- [ ] `configure` — elementary overlay from `preferences-system`
  names: configure · category: actions
  16 ← categories/16/preferences-system.svg
  24, 22 ← categories/24/preferences-system.svg
  32 ← categories/32/preferences-system.svg
  48 ← categories/48/preferences-system.svg
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  note: Settings → Configure <App> (KStandardAction::Preferences). elementary's settings gear.
  note: elementary art at 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: albumfolder-properties (16, 22, 24), cm_options (16, 22, 24), configuration (16, 22, 24), configure-symbolic (22, 32), configure_kile (16, 22, 24), configure_project (16, 22, 24), cvc-configure (16, 22), dialog-object-properties (16, 22, 24), games-config-options (16, 22, 24), gtk-execute (16, 22, 24), gtk-properties (16, 22, 24), hb-presets (16, 22, 24), itmages-settings (16, 22, 24), preferences-desktop (16), preferences-other (16, 22, 24), rabbitvcs-settings (16, 22, 24), services (16, 22, 24), settings (16, 22, 24), step_object_Controller (16, 22, 24), stock_file-properites (16, 22, 24), stock_folder_properties (16, 22, 24), stock_properties (16, 22, 24), user-group-properties (16, 22, 24), user-properties (16, 22, 24), vm-settings (16, 22, 24), vmware-settings (16, 22, 24), xfce-system-settings (16, 22, 24), xfce4-settings (16, 22, 24)
- [ ] `system-run` — elementary overlay
  names: system-run, system-run-symbolic · category: actions
  16 ← actions/16/system-run.svg
  24, 22 ← actions/24/system-run.svg
  32 ← actions/32/system-run.svg
  48 ← actions/48/system-run.svg
  system-run-symbolic: symbolic ← actions/symbolic/system-run-symbolic.svg (recolor to ColorScheme-Text)
  note: elementary art at 16, 24, 32, 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: CVnamespace (16, 22, 24), featured-vm (16, 22, 24), gnome-run (16), gtk-preferences (16, 22, 24), package_settings (16, 22, 24), preferences-system-session-services (16, 22, 24), rabbitvcs-dbus (16, 22, 24), view-process-system (16, 22, 24)
- [ ] `help-about` — elementary overlay
  names: help-about, help-about-symbolic · category: actions
  16 ← actions/16/help-about.svg
  24, 22 ← actions/24/help-about.svg
  32 ← actions/32/help-about.svg
  48 ← actions/48/help-about.svg
  help-about-symbolic: symbolic ← actions/symbolic/help-about-symbolic.svg (recolor to ColorScheme-Text)
  also changes: cvc-info (16, 22), exifinfo (16, 22, 24), gnome-app-install-star (16, 22, 24), gtk-about (16, 22, 24), gtk-about-symbolic (symbolic), help-feedback (16, 22, 24), help-hint (16, 22, 24), help-latex (16, 22, 24), hook-notifier (16, 22, 24), icon_details (16, 22, 24), kt-info-widget (16, 22, 24), rabbitvcs-about (16, 22, 24), setroubleshoot_icon (16, 22, 24), stock_about (16, 22, 24), stock_navigator (16, 22, 24), xsi-help-about-symbolic (symbolic)
- [ ] `tab-close` — elementary overlay from `window-close`
  names: tab-close · category: actions
  16 ← actions/16/window-close.svg
  24, 22 ← actions/24/window-close.svg
  48 ← actions/48/window-close.svg
  note: Close Tab (Dolphin, Konsole, Kate, Okular).
  note: elementary art at 16, 24, 48px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: tab-close-other (16, 22, 24), tab_remove (16, 22, 24), tab_remove_other (16, 22, 24)
- [ ] `configure-shortcuts` — elementary overlay from `preferences-desktop-keyboard`
  names: configure-shortcuts · category: actions
  16 ← categories/16/preferences-desktop-keyboard.svg
  24, 22 ← categories/24/preferences-desktop-keyboard.svg
  32 ← categories/32/preferences-desktop-keyboard.svg
  48 ← categories/48/preferences-desktop-keyboard.svg
  note: KDE standard action: referenced by KDE's standard-actions library (libKF6ConfigWidgets/libKF6ConfigGui), so it likely shows in the File/Edit/View/Settings/Help menus of every KDE app
  note: Settings → Configure Keyboard Shortcuts.
  note: elementary art at 16, 24, 32px is grayscale — it will not follow the KDE color scheme, check it on the dark preview
  also changes: configure_shortcuts (16, 22, 24), fcitx-kbd (16, 22, 24), help-keybord-shortcuts (16, 22, 24), view-media-config (22, 24)
- [ ] `notification-disabled` — elementary overlay
  names: notification-disabled, notification-disabled-symbolic · category: actions
  48, 16, 22, 24 ← status/48/notification-disabled.svg
  notification-disabled-symbolic: symbolic ← status/symbolic/notification-disabled-symbolic.svg (recolor to ColorScheme-Text)
