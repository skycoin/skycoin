This folder contains the GUI translation files. To maintain order and be able
to easily make any necessary updates to the translation files after updating
the main text file, please follow its instructions if you are working with
its contents.

# Contents of this folder

The contents of this folder are:

- `README.md`: this file.

- `check.js`: file with the script for detecting if a translation file has errors
or should be updated.

- `en.json`: main file with all the texts of the application, in English. It should
only be modified when changing the texts of the application (add, modify and
delete). This means that the file must not be modified while creating a new 
ranslation or modifying an existing one.

- Various `xx.json` files: files with the translated versions of the texts of
`en.json`.

- Various `xx_base.json` files: files with copies of `en.json` made the last time the
corresponding `xx.json` file was modified.

Normally there is no need to modify the first two files.

For more information about the `xx.json` and `xx_base.json`, please check the
[Add a new translation](#add-a-new-translation) and
[Update a translation](#update-a-translation) sections.

# About the meaning of "xx" in this file

Several parts of this file uses "xx" as part of file names or scripts, like
`xx.json` and `xx_base.json`. In fact, no file in this folder should be called
`xx.json` or `xx_base.json`, the "xx" part must be replaces with the two
characters code of the languaje. For example, if you are working with the Chinese
translation, the files will be `cn.json` and `cn_base.json`, instead of `xx.json`
and `xx_base.json`. The same if true for the scripts, if you are working with the
Chinese translation, instead of running `node check.js xx` you must run
`node check.js cn`.

# Add a new translation

First you must create in this folder two copies of the `en.json` file. The first
copy must be called `xx.json`, where the `xx` part must be the two characters code
of the new languaje. For example, for Chinese the name of the file should be
`cn.json`; for Spanish, `es.json`; for French, `fr.json`, etc.

The second copy of `en.json` must be renamed to `xx_base.json`, where the `xx` part
must be the two characters code of the new languaje. This means that if the first
copy is named `cn.json`, the second one should be named `cn_base.json`.

It is not necessary to follow a specific standard for the two characters code, but
it must be limited to two letters and be a recognizable code for the language.

After creating the two files, simply translate the texts in `xx.json`. Please make
sure you do not modify the structure of `xx.json`, just modify the texts.

The `xx_base.json` file must not be modified in any way, as it is used only as a way
to know what the state of `en.json` was the last time the `xx.json` file was
modified. This copy will be compared in the future with `en.json`, to verify if
there were modifications to `en.json` since the last time the translation file was
modified and if an update is needed.

If the `xx.json` and `xx_base.json` files do not have the same elements, the
automatic tests could fail when uploading the changes to the repository, preventing
the changes from being accepted, so, again, it is important not to modify the
structure of `xx.json`, but only its contents.

After doing all this, the translation will be ready, but will not be available in
the GUI until adding it to the code.

# Verify the translation files

This folder includes a script that is capable of automatically checking the
translation files, to detect problems and know what should be updated.

For using it, your computer must have `Node.js` installed.

## Checking for problems

For detecting basic problems on the translation files, open a command line window
in this folder and run `node check.js`. This will check the following:

- The `en.json` must exist, as it is the main languaje file for the app.

- For every `xx.json` file (except `en.json`) an `xx_base.json` file must exist
and viceversa.

- A `xx.json` file and its corresponding `xx_base.json` file must have the exact
same elements (only the content of that elements could be different), as the
`xx.json` is suposed to be the translation of the contents of `xx_base.json`.

As you can see, this only checks for errors that could be made while creating or
modifying the `xx.json` and `xx_base.json` files, and does not check if any
translation needs to be updated.

At the end of the script excecution, the console will display the list of all
errors found, if any. This check could be done automatically when making changes
to the repository, to reject updates with problems, so it is good idea to run it
manually before uploading changes.

Note: at this time the script does not check if the elements of the files are
in the same order, but this could be added in the future, so it is recomended
not to change the order of the elements.

## Checking if a languaje file needs to be updated

To detect if an specific languaje needs updating, run `node check.js xx`,
where xx is the two characters code of the languaje you want to check. If you
want to check all languajes, run `node check.js all`.

By doing this, the script will perform all the checks described in the
[Checking for problems](#checking-for-problems) section, plus this:

- The `en.json` and `xx_base.json` should have the same elements. If `en.json`
has elements that `xx_base.json` does not contain, it means that, since the
last time the translation file was updated, new texts have been added to the
application. If `xx_base.json` has elements that `en.json` does not contain,
it means that, since the last time the translation file was updated, some texts
have been removed from the application. Both cases mean that the translation
file should be updated.

- The elements of `en.json` and `xx_base.json` should have the same content.
If any element have different content, it means that since the last time the
translation file was updated, some texts of the applications have been changed.
This means that the translation file should be updated.

At the end of the script excecution, the console will display the list of all
errors found, if any.

# Update a translation

Before updating a translation file, you should follow the steps of the
[Checking if a languaje file needs to be updated](#Checking-if-a-languaje-file-needs-to-be-updated)
section. By doing so you will quikly know exactly what texts must be added,
deleted or edited.

After doing that, make all the required modifications in the `xx.json` file,
this menans adding, deleting and modifying all the elements indicated by the
script. Please be sure to modify only what is required and to add any new
element in the same position that it is in the `en.json` file. This process
is manual, so be sure check all the changes before finishing.

After doing the modifications in `xx.json`, delete the `xx_base.json` file,
create a copy of `en.json` and rename it `xx_base.json`. The objetive is to
simply update the `xx_base.json` file to the current state of `en.json`.
this will make possible to check in the future if more updates are nedded,
due to new changes in `en.json`.

Once all the changes are made, check again the languaje file as indicated
in the
[Checking if a languaje file needs to be updated](#Checking-if-a-languaje-file-needs-to-be-updated)
section. The script should not return errors. If the script returns errors,
please solve them before continuing.
