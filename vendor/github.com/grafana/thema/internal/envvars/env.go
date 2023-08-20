package envvars

import "os"

// VarUpdateGolden is the name of the env var to trigger updating golden test files.
const VarUpdateGolden = "THEMA_UPDATE_GOLDEN"

// ForceVerify indicates that all verifications should be performed, even if
// e.g. SkipBuggyChecks() says otherwise.
var ForceVerify = os.Getenv("THEMA_FORCEVERIFY") != ""

// ReverseTranslate indicates whether reverse translation is supported.
//
// Used primarily as a single point of control for testing.
//
// Permanently set to true, as reverse translation is now supported.
var ReverseTranslate = true

// var ReverseTranslate = os.Getenv("THEMA_REVERSETRANSLATE") != ""

// UpdateGoldenFiles determines whether testscript scripts should update txtar
// archives in the event of cmp failures.
// It is controlled by setting THEMA_UPDATE_GOLDEN to a non-empty string like "true".
// It corresponds to testscript.Params.UpdateGoldenFiles; see its docs for details.
var UpdateGoldenFiles = os.Getenv(VarUpdateGolden) != ""

// FormatTxtar ensures that .cue files in txtar test archives are well
// formatted, updating the archive as required prior to running a test.
// It is controlled by setting THEMA_FORMAT_TXTAR to a non-empty string like "true".
var FormatTxtar = os.Getenv("THEMA_FORMAT_TXTAR") != ""

// FixLineages governs whether lineages defined in txtar files should be fixed
// and rewritten on the fly.
var FixLineages = os.Getenv("THEMA_FIX_TXTAR_LINEAGES") != ""
