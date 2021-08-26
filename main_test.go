// Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseDiff(t *testing.T) {
	testName := "Multiple files"
	diff := []byte(`
diff --git a/README.md b/README.md
index d4edde0..807b76d 100644
--- a/README.md
+++ b/README.md
@@ -1,0 +2 @@
+hello
diff --git a/repositories.txt b/repositories.txt
index cff484d..e14c179 100644
--- a/repositories.txt
+++ b/repositories.txt
@@ -8,0 +9 @@ https://github.com/arduino-libraries/Ethernet
+https://github.com/foo/bar
`)

	requestType, requestError, arduinoLintLibraryManagerSetting, submissionURLs := parseDiff(diff, "repositories.txt")
	assert.Equal(t, "other", requestType, testName)
	assert.Equal(t, "", requestError, testName)
	assert.Equal(t, "", arduinoLintLibraryManagerSetting, testName)
	assert.Nil(t, submissionURLs, testName)

	testName = "Not list"
	diff = []byte(`
diff --git a/README.md b/README.md
index d4edde0..807b76d 100644
--- a/README.md
+++ b/README.md
@@ -1 +1,2 @@
	# Arduino Library Manager list
+hello
`)

	requestType, requestError, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "other", requestType, testName)
	assert.Equal(t, "", requestError, testName)
	assert.Equal(t, "", arduinoLintLibraryManagerSetting, testName)
	assert.Nil(t, submissionURLs, testName)

	testName = "List filename change"
	diff = []byte(`
diff --git a/repositories.txt b/foobar.txt
similarity index 99%
rename from repositories.txt
rename to foobar.txt
index cff484d..e14c179 100644
--- a/repositories.txt
+++ b/foobar.txt
@@ -8,0 +9 @@ https://github.com/arduino-libraries/Ethernet
+https://github.com/foo/bar
`)

	requestType, requestError, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "other", requestType, testName)
	assert.Equal(t, "", requestError, testName)
	assert.Equal(t, "", arduinoLintLibraryManagerSetting, testName)
	assert.Nil(t, submissionURLs, testName)

	testName = "Submission"
	diff = []byte(`
diff --git a/repositories.txt b/repositories.txt
index cff484d..9f67763 100644
--- a/repositories.txt
+++ b/repositories.txt
@@ -8,0 +9,2 @@ https://github.com/arduino-libraries/Ethernet
+https://github.com/foo/bar
+https://github.com/foo/baz
`)

	requestType, requestError, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "submission", requestType, testName)
	assert.Equal(t, "", requestError, testName)
	assert.Equal(t, "submit", arduinoLintLibraryManagerSetting, testName)
	assert.ElementsMatch(t, []string{"https://github.com/foo/bar", "https://github.com/foo/baz"}, submissionURLs, testName)

	testName = "Submission w/ no newline at end of file"
	diff = []byte(`
diff --git a/repositories.txt b/repositories.txt
index cff484d..1b0b80b 100644
--- a/repositories.txt
+++ b/repositories.txt
@@ -3391,0 +3392 @@ https://github.com/lbernstone/plotutils
+https://github.com/foo/bar
\ No newline at end of file
`)

	requestType, requestError, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "invalid", requestType, testName)
	assert.Equal(t, "Pull request removes newline from the end of a file.%0APlease add a blank line to the end of the file.", requestError, testName)
	assert.Equal(t, "", arduinoLintLibraryManagerSetting, testName)
	assert.Nil(t, submissionURLs, testName)

	testName = "Submission w/ blank line"
	diff = []byte(`
diff --git a/repositories.txt b/repositories.txt
index cff484d..1b0b80b 100644
--- a/repositories.txt
+++ b/repositories.txt
@@ -3391,0 +3392 @@ https://github.com/lbernstone/plotutils
+https://github.com/foo/bar
+
`)

	requestType, requestError, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "submission", requestType, testName)
	assert.Equal(t, "", requestError, testName)
	assert.Equal(t, "submit", arduinoLintLibraryManagerSetting, testName)
	assert.ElementsMatch(t, []string{"https://github.com/foo/bar"}, submissionURLs, testName)

	testName = "Removal"
	diff = []byte(`
diff --git a/repositories.txt b/repositories.txt
index cff484d..38e11d8 100644
--- a/repositories.txt
+++ b/repositories.txt
@@ -8 +7,0 @@ https://github.com/firmata/arduino
-https://github.com/arduino-libraries/Ethernet
`)

	requestType, requestError, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "removal", requestType, testName)
	assert.Equal(t, "", requestError, testName)
	assert.Equal(t, "", arduinoLintLibraryManagerSetting, testName)
	assert.Nil(t, submissionURLs, testName)

	testName = "Modification"
	diff = []byte(`
diff --git a/repositories.txt b/repositories.txt
index cff484d..8b401a1 100644
--- a/repositories.txt
+++ b/repositories.txt
@@ -8 +8 @@ https://github.com/firmata/arduino
-https://github.com/arduino-libraries/Ethernet
+https://github.com/foo/bar
`)

	requestType, requestError, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "modification", requestType, testName)
	assert.Equal(t, "", requestError, testName)
	assert.Equal(t, "update", arduinoLintLibraryManagerSetting, testName)
	assert.Equal(t, []string{"https://github.com/foo/bar"}, submissionURLs, testName)

	testName = "Newline-only"
	diff = []byte(`
diff --git a/repositories.txt b/repositories.txt
index d9a6136..ca902d9 100644
--- a/repositories.txt
+++ b/repositories.txt
@@ -1,4 +1,4 @@
	https://github.com/firmata/arduino
-
	https://github.com/arduino-libraries/Ethernet
	https://github.com/lbernstone/plotutils
+
`)

	requestType, requestError, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "other", requestType, testName)
	assert.Equal(t, "", requestError, testName)
	assert.Equal(t, "", arduinoLintLibraryManagerSetting, testName)
	assert.Nil(t, submissionURLs, testName)
}

func Test_normalizeURL(t *testing.T) {
	testTables := []struct {
		testName              string
		rawURL                string
		expectedNormalizedURL string
	}{
		{"Trailing slash", "https://github.com/foo/bar/", "https://github.com/foo/bar.git"},
		{".git suffix", "https://github.com/foo/bar.git", "https://github.com/foo/bar.git"},
		{"http://", "http://github.com/foo/bar", "https://github.com/foo/bar.git"},
		{"git://", "git://github.com/foo/bar", "https://github.com/foo/bar.git"},
		{"Root URL", "https://github.com", "https://github.com/"},
		{"Root URL with trailing slash", "https://github.com/", "https://github.com/"},
	}

	for _, testTable := range testTables {
		rawURL, err := url.Parse(testTable.rawURL)
		require.Nil(t, err)
		expectedNormalizedURL, err := url.Parse(testTable.expectedNormalizedURL)
		require.Nil(t, err)

		assert.Equal(t, *expectedNormalizedURL, normalizeURL(rawURL), testTable.testName)
	}
}

func Test_indexerLogsURL(t *testing.T) {
	testTables := []struct {
		testName               string
		normalizedURL          string
		expectedIndexerLogsURL string
	}{
		{"GitHub", "https://github.com/foo/bar.git", "http://downloads.arduino.cc/libraries/logs/github.com/foo/bar/"},
		{"GitLab", "https://gitlab.com/yesbotics/libs/arduino/voltmeter.git", "http://downloads.arduino.cc/libraries/logs/gitlab.com/yesbotics/libs/arduino/voltmeter/"},
	}

	for _, testTable := range testTables {
		assert.Equal(t, testTable.expectedIndexerLogsURL, indexerLogsURL(testTable.normalizedURL), testTable.testName)
	}
}

func Test_uRLIsUnder(t *testing.T) {
	testTables := []struct {
		testName         string
		childURL         string
		parentCandidates []string
		assertion        assert.BoolAssertionFunc
	}{
		{"Match, root path", "https://github.com/foo/bar", []string{"example.com", "github.com"}, assert.True},
		{"Mismatch, root path", "https://github.com/foo/bar", []string{"example.com", "example.org"}, assert.False},
		{"Match, subfolder", "https://github.com/foo/bar", []string{"example.com/foo", "github.com/foo"}, assert.True},
		{"Mismatch, subfolder", "https://github.com/foo/bar", []string{"example.com/foo", "github.org/bar"}, assert.False},
		{"Match, root child URL", "https://github.com/", []string{"example.com", "github.com"}, assert.True},
		{"Mismatch, root child URL", "https://github.com/", []string{"example.com", "github.org"}, assert.False},
	}

	for _, testTable := range testTables {
		childURL, err := url.Parse(testTable.childURL)
		require.Nil(t, err)

		t.Run(testTable.testName, func(t *testing.T) {
			testTable.assertion(t, uRLIsUnder(*childURL, testTable.parentCandidates))
		})
	}
}
