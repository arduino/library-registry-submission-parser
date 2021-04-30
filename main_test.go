// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
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

	requestType, arduinoLintLibraryManagerSetting, submissionURLs := parseDiff(diff, "repositories.txt")
	assert.Equal(t, "other", requestType, testName)
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

	requestType, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "other", requestType, testName)
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

	requestType, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "other", requestType, testName)
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

	requestType, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "submission", requestType, testName)
	assert.Equal(t, "submit", arduinoLintLibraryManagerSetting, testName)
	assert.ElementsMatch(t, submissionURLs, []string{"https://github.com/foo/bar", "https://github.com/foo/baz"}, testName)

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

	requestType, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "submission", requestType, testName)
	assert.Equal(t, "submit", arduinoLintLibraryManagerSetting, testName)
	assert.ElementsMatch(t, submissionURLs, []string{"https://github.com/foo/bar"}, testName)

	testName = "Submission w/ blank line"
	diff = []byte(`
diff --git a/repositories.txt b/repositories.txt
index cff484d..1b0b80b 100644
--- a/repositories.txt
+++ b/repositories.txt
@@ -3391,0 +3392 @@ https://github.com/lbernstone/plotutils
+https://github.com/foo/bar
\ No newline at end of file
`)

	requestType, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "submission", requestType, testName)
	assert.Equal(t, "submit", arduinoLintLibraryManagerSetting, testName)
	assert.ElementsMatch(t, submissionURLs, []string{"https://github.com/foo/bar"}, testName)

	testName = "Removal"
	diff = []byte(`
diff --git a/repositories.txt b/repositories.txt
index cff484d..38e11d8 100644
--- a/repositories.txt
+++ b/repositories.txt
@@ -8 +7,0 @@ https://github.com/firmata/arduino
-https://github.com/arduino-libraries/Ethernet
`)

	requestType, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "removal", requestType, testName)
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

	requestType, arduinoLintLibraryManagerSetting, submissionURLs = parseDiff(diff, "repositories.txt")
	assert.Equal(t, "modification", requestType, testName)
	assert.Equal(t, "update", arduinoLintLibraryManagerSetting, testName)
	assert.Equal(t, submissionURLs, []string{"https://github.com/foo/bar"}, testName)
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
