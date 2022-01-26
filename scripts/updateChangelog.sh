#!/usr/bin/env bash
# Credits to github.com/particleflux for creating this in a different project

newVersion="$1"
currentVersion=$(grep -E '^##\s+\[[0-9.]+\]' CHANGELOG.md | head -n1 | awk '{print $2}' | tr -d [])
echo "Current version is $currentVersion"
echo "New version is $newVersion"

# replace "## [Unreleased]" with new version headline "## [x.y.z] - date"
newHeadline="\#\# \[$newVersion\] - $(date -I)"
sed -i "s/\#\# \[Unreleased\]/$newHeadline/i" CHANGELOG.md
# re-prepend [Unreleased]
sed -i -e "/^$newHeadline\s*$/i ## [Unreleased]\n\n" CHANGELOG.md

# Update links at the end
sed -i -e "s#v${currentVersion}...HEAD#v$newVersion...HEAD#i" CHANGELOG.md
sed -i -e "/^\[Unreleased\]: /a      [$newVersion]: https://github.com/DBX12/envManager/compare/v${currentVersion}...v${newVersion}" CHANGELOG.md