Release process
===============

This document outlines how to create a release of net/metrics

1.  Set up some environment variables for use later.

    ```
    # This is the version being released.
    VERSION=1.21.0
    ```

2.  Make sure you have the latest master.

    ```
    git checkout master
    git pull
    ```

3.  Alter the release date in CHANGELOG.md for `$VERSION` using the format
    `YYYY-MM-DD` and remove the trailing `-dev`, making the latest version
    match `$VERSION`.

    ```diff
    -## v1.21.0-dev (unreleased)
    +## v1.21.0 (2017-10-23)
    ```

4.  Update the version number in version.go and verify that it matches what is
    in the changelog.

    ```
    sed -i '' -e "s/^const Version =.*/const Version = \"$VERSION\"/" version.go
    make verifyversion
    ```

5.  Create a commit for the release.

    ```
    git add version.go CHANGELOG.md
    git commit -m "Preparing release v$VERSION"
    ```

6.  Tag and push the release.

    ```
    git tag -a "v$VERSION" -m "v$VERSION"
    git push origin master "v$VERSION"
    ```

7.  Go to <https://github.com/yarpc/metrics/tags> and edit the release notes
    of the new tag.  Copy the changelog entries for this release int the
    release notes and set the name of the release to the version number
    (`v$VERSION`).

8.  Add a placeholder for the next version to CHANGELOG.md.  This is typically
    one minor version above the version just released.

    ```diff
    +v1.22.0-dev (unreleased)
    +------------------------
    +
    +-   No changes yet.
    +
    +
     v1.21.0 (2017-10-23)
     --------------------
    ```

9.  Update the version number in version.go to the same version.

    ```diff
    -const Version = "1.21.0"
    +const Version = "1.22.0-dev"
    ```

10. Verify the version number matches.

    ```
    make verifyversion
    ```

11. Commit and push your changes.

    ```
    git add CHANGELOG.md version.go
    git commit -m 'Back to development'
    git push origin master
    ```
