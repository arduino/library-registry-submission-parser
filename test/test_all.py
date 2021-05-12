# Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
#
# This software is released under the GNU General Public License version 3.
# The terms of this license can be found at:
# https://www.gnu.org/licenses/gpl-3.0.en.html
#
# You can be released from the requirements of the above licenses by purchasing
# a commercial license. Buying such a license is mandatory if you want to
# modify or otherwise use the software for commercial activities involving the
# Arduino software without disclosing the source code of your own applications.
# To purchase a commercial license, send an email to license@arduino.cc.
import json
import pathlib
import platform
import typing

import invoke.context
import pytest

test_data_path = pathlib.Path(__file__).resolve().parent.joinpath("testdata")


@pytest.mark.parametrize(
    "repopath_folder_name, expected_type, expected_submissions, expected_indexentry, expected_indexerlogsurls",
    [
        ("list-deleted-diff", "other", None, "", ""),
        ("multi-file-diff", "other", None, "", ""),
        ("non-list-diff", "other", None, "", ""),
        ("list-rename-diff", "other", None, "", ""),
        ("removal", "removal", None, "", ""),
        (
            "modification",
            "modification",
            [
                {
                    "submissionURL": "https://github.com/arduino-libraries/ArduinoCloudThing",
                    "normalizedURL": "https://github.com/arduino-libraries/ArduinoCloudThing.git",
                    "repositoryName": "ArduinoCloudThing",
                    "name": "ArduinoCloudThing",
                    "official": True,
                    "tag": "1.7.3",
                    "error": "",
                }
            ],
            "https://github.com/arduino-libraries/ArduinoCloudThing.git|Arduino|ArduinoCloudThing",
            "http://downloads.arduino.cc/libraries/logs/github.com/arduino-libraries/ArduinoCloudThing/",
        ),
        (
            "url-error",
            "submission",
            [
                {
                    "submissionURL": "foo",
                    "normalizedURL": "",
                    "repositoryName": "",
                    "name": "",
                    "official": False,
                    "tag": "",
                    "error": 'Unable to load submission URL: Get "foo": unsupported protocol scheme ""',
                }
            ],
            "",
            "http://downloads.arduino.cc/libraries/logs//",
        ),
        (
            "url-404",
            "submission",
            [
                {
                    "submissionURL": "http://httpstat.us/404",
                    "normalizedURL": "",
                    "repositoryName": "",
                    "name": "",
                    "official": False,
                    "tag": "",
                    "error": "Unable to load submission URL. Is the repository public?",
                }
            ],
            "",
            "http://downloads.arduino.cc/libraries/logs//",
        ),
        (
            "not-supported-git-host",
            "submission",
            [
                {
                    "submissionURL": "https://example.com",
                    "normalizedURL": "https://example.com/",
                    "repositoryName": "",
                    "name": "",
                    "official": False,
                    "tag": "",
                    "error": "example.com is not currently supported as a Git hosting website for Library Manager.%0A"
                    "See: https://github.com/arduino/library-registry/blob/main/FAQ.md#what-are-the-requirements-for-a-"
                    "library-to-be-added-to-library-manager",
                }
            ],
            "",
            "http://downloads.arduino.cc/libraries/logs/example.com//",
        ),
        (
            "not-git-clone-url",
            "submission",
            [
                {
                    "submissionURL": "https://github.com/arduino-libraries/ArduinoCloudThing/releases",
                    "normalizedURL": "https://github.com/arduino-libraries/ArduinoCloudThing/releases.git",
                    "repositoryName": "",
                    "name": "",
                    "official": False,
                    "tag": "",
                    "error": "Submission URL is not a Git clone URL (e.g., `https://github.com/arduino-libraries/Servo`"
                    ").",
                }
            ],
            "",
            "http://downloads.arduino.cc/libraries/logs/github.com/arduino-libraries/ArduinoCloudThing/releases/",
        ),
        (
            "already-in-library-manager",
            "submission",
            [
                {
                    "submissionURL": "https://github.com/arduino-libraries/Servo",
                    "normalizedURL": "https://github.com/arduino-libraries/Servo.git",
                    "repositoryName": "Servo",
                    "name": "",
                    "official": False,
                    "tag": "",
                    "error": "Submission URL is already in the Library Manager index.",
                }
            ],
            "",
            "http://downloads.arduino.cc/libraries/logs/github.com/arduino-libraries/Servo/",
        ),
        (
            "type-arduino",
            "submission",
            [
                {
                    "submissionURL": "https://github.com/arduino-libraries/ArduinoCloudThing",
                    "normalizedURL": "https://github.com/arduino-libraries/ArduinoCloudThing.git",
                    "repositoryName": "ArduinoCloudThing",
                    "name": "ArduinoCloudThing",
                    "official": True,
                    "tag": "1.7.3",
                    "error": "",
                }
            ],
            "https://github.com/arduino-libraries/ArduinoCloudThing.git|Arduino|ArduinoCloudThing",
            "http://downloads.arduino.cc/libraries/logs/github.com/arduino-libraries/ArduinoCloudThing/",
        ),
        (
            "type-partner",
            "submission",
            [
                {
                    "submissionURL": "https://github.com/ms-iot/virtual-shields-arduino",
                    "normalizedURL": "https://github.com/ms-iot/virtual-shields-arduino.git",
                    "repositoryName": "virtual-shields-arduino",
                    "name": "Windows Virtual Shields for Arduino",
                    "official": False,
                    "tag": "v1.2.0",
                    "error": "",
                }
            ],
            "https://github.com/ms-iot/virtual-shields-arduino.git|Partner|Windows Virtual Shields for Arduino",
            "http://downloads.arduino.cc/libraries/logs/github.com/ms-iot/virtual-shields-arduino/",
        ),
        (
            "type-recommended",
            "submission",
            [
                {
                    "submissionURL": "https://github.com/adafruit/Adafruit_TinyFlash",
                    "normalizedURL": "https://github.com/adafruit/Adafruit_TinyFlash.git",
                    "repositoryName": "Adafruit_TinyFlash",
                    "name": "Adafruit TinyFlash",
                    "official": False,
                    "tag": "1.0.4",
                    "error": "",
                }
            ],
            "https://github.com/adafruit/Adafruit_TinyFlash.git|Recommended|Adafruit TinyFlash",
            "http://downloads.arduino.cc/libraries/logs/github.com/adafruit/Adafruit_TinyFlash/",
        ),
        (
            "type-contributed",
            "submission",
            [
                {
                    "submissionURL": "https://github.com/sparkfun/SparkFun_Ublox_Arduino_Library",
                    "normalizedURL": "https://github.com/sparkfun/SparkFun_Ublox_Arduino_Library.git",
                    "repositoryName": "SparkFun_Ublox_Arduino_Library",
                    "name": "SparkFun u-blox Arduino Library",
                    "official": False,
                    "tag": "v1.8.11",
                    "error": "",
                }
            ],
            "https://github.com/sparkfun/SparkFun_Ublox_Arduino_Library.git|Contributed|SparkFun u-blox Arduino Library"
            "",
            "http://downloads.arduino.cc/libraries/logs/github.com/sparkfun/SparkFun_Ublox_Arduino_Library/",
        ),
        (
            "no-tags",
            "submission",
            [
                {
                    "submissionURL": "https://github.com/arduino/cloud-examples",
                    "normalizedURL": "https://github.com/arduino/cloud-examples.git",
                    "repositoryName": "cloud-examples",
                    "name": "",
                    "official": True,
                    "tag": "",
                    "error": "The repository has no tags. You need to create a [release](https://docs.github.com/en/git"
                    "hub/administering-a-repository/managing-releases-in-a-repository) or [tag](https://git-scm.com/doc"
                    "s/git-tag) that matches the `version` value in the library's library.properties file.",
                }
            ],
            "",
            "http://downloads.arduino.cc/libraries/logs/github.com/arduino/cloud-examples/",
        ),
        (
            "no-library-properties",
            "submission",
            [
                {
                    "submissionURL": "https://github.com/arduino-libraries/WiFiLink-Firmware",
                    "normalizedURL": "https://github.com/arduino-libraries/WiFiLink-Firmware.git",
                    "repositoryName": "WiFiLink-Firmware",
                    "name": "",
                    "official": True,
                    "tag": "1.0.1",
                    "error": "Library is missing a library.properties metadata file.",
                }
            ],
            "",
            "http://downloads.arduino.cc/libraries/logs/github.com/arduino-libraries/WiFiLink-Firmware/",
        ),
        (
            "duplicates-in-submission",
            "submission",
            [
                {
                    "submissionURL": "https://github.com/arduino-libraries/ArduinoCloudThing",
                    "normalizedURL": "https://github.com/arduino-libraries/ArduinoCloudThing.git",
                    "repositoryName": "ArduinoCloudThing",
                    "name": "ArduinoCloudThing",
                    "official": True,
                    "tag": "1.7.3",
                    "error": "",
                },
                {
                    "submissionURL": "https://github.com/arduino-libraries/ArduinoCloudThing",
                    "normalizedURL": "https://github.com/arduino-libraries/ArduinoCloudThing.git",
                    "repositoryName": "ArduinoCloudThing",
                    "name": "ArduinoCloudThing",
                    "official": True,
                    "tag": "1.7.3",
                    "error": "Submission contains duplicate URLs.",
                },
            ],
            "https://github.com/arduino-libraries/ArduinoCloudThing.git|Arduino|ArduinoCloudThing%0A"
            "https://github.com/arduino-libraries/ArduinoCloudThing.git|Arduino|ArduinoCloudThing",
            "http://downloads.arduino.cc/libraries/logs/github.com/arduino-libraries/ArduinoCloudThing/%0A"
            "http://downloads.arduino.cc/libraries/logs/github.com/arduino-libraries/ArduinoCloudThing/",
        ),
    ],
)
def test_request(
    run_command,
    repopath_folder_name,
    expected_type,
    expected_submissions,
    expected_indexentry,
    expected_indexerlogsurls,
):
    diffpath = test_data_path.joinpath(repopath_folder_name, "diff.txt")
    repopath = test_data_path.joinpath(repopath_folder_name)
    listname = "repositories.txt"

    result = run_command(cmd=["--diffpath", diffpath, "--repopath", repopath, "--listname", listname])
    assert result.ok

    request = json.loads(result.stdout)
    assert request["type"] == expected_type
    assert request["submissions"] == expected_submissions
    assert request["indexEntry"] == expected_indexentry
    assert request["submissions"] == expected_submissions
    assert request["indexerLogsURLs"] == expected_indexerlogsurls


@pytest.fixture(scope="function")
def run_command(pytestconfig, working_dir) -> typing.Callable[..., invoke.runners.Result]:
    """Provide a wrapper around invoke's `run` API so that every test will work in the same temporary folder.

    Useful reference:
        http://docs.pyinvoke.org/en/1.4/api/runners.html#invoke.runners.Result
    """

    executable_path = pathlib.Path(pytestconfig.rootdir).parent / "parser"

    def _run(
        cmd: list, custom_working_dir: typing.Optional[str] = None, custom_env: typing.Optional[dict] = None
    ) -> invoke.runners.Result:
        if cmd is None:
            cmd = []
        if not custom_working_dir:
            custom_working_dir = working_dir
        quoted_cmd = []
        for token in cmd:
            quoted_cmd.append(f'"{token}"')
        cli_full_line = '"{}" {}'.format(executable_path, " ".join(quoted_cmd))
        run_context = invoke.context.Context()
        # It might happen that we need to change directories between drives on Windows,
        # in that case the "/d" flag must be used otherwise directory wouldn't change
        cd_command = "cd"
        if platform.system() == "Windows":
            cd_command += " /d"
        # Context.cd() is not used since it doesn't work correctly on Windows.
        # It escapes spaces in the path using "\ " but it doesn't always work,
        # wrapping the path in quotation marks is the safest approach
        with run_context.prefix(f'{cd_command} "{custom_working_dir}"'):
            return run_context.run(
                command=cli_full_line, echo=False, hide=True, warn=True, env=custom_env, encoding="utf-8"
            )

    return _run


@pytest.fixture(scope="function")
def working_dir(tmpdir_factory) -> str:
    """Create a temporary folder for the test to run in. It will be created before running each test and deleted at the
    end. This way all the tests work in isolation.
    """
    work_dir = tmpdir_factory.mktemp(basename="TestWorkingDir")
    yield str(work_dir)
