// Package internal
package internal

import "github.com/spyzhov/ajson/v1/internal"

type (
	State = internal.State
	Class = internal.Class

	States [54][38]State
)

const __ = -1
const ѢѢ = -1 // fixme: not set yet

// enum classes
const (
	C_SPACE Class = iota /* space */
	C_WHITE              /* other whitespace */
	C_LCURB              /* {  */
	C_RCURB              /* } */
	C_LSQRB              /* [ */
	C_RSQRB              /* ] */
	C_COLON              /* : */
	C_COMMA              /* , */
	C_QUOTE              /* "|' */
	C_BACKS              /* \ */
	C_SLASH              /* / */
	C_PLUS               /* + */
	C_MINUS              /* - */
	C_POINT              /* . */
	C_ZERO               /* 0 */
	C_DIGIT              /* 123456789 */
	C_LOW_A              /* a */
	C_LOW_B              /* b */
	C_LOW_C              /* c */
	C_LOW_D              /* d */
	C_LOW_E              /* e */
	C_LOW_F              /* f */
	C_LOW_L              /* l */
	C_LOW_N              /* n */
	C_LOW_R              /* r */
	C_LOW_S              /* s */
	C_LOW_T              /* t */
	C_LOW_U              /* u */
	C_ABCDF              /* ABCDF */
	C_E                  /* E */
	C_LPARB              /* ( */
	C_RPARB              /* ) */
	C_QUESM              /* ? */
	C_DOLAR              /* $ */
	C_AT                 /* @ */
	C_UNDER              /* _ */
	C_ASTER              /* * */
	C_ETC                /* everything else */
)

// AsciiClasses array maps the 128 ASCII characters into character classes.
var AsciiClasses = [128]Class{
	/*
	   This array maps the 128 ASCII characters into character classes.
	   The remaining Unicode characters should be mapped to C_ETC.
	   Non-whitespace control characters are errors.
	*/
	__, __, __, __, __, __, __, __,
	__, C_WHITE, C_WHITE, __, __, C_WHITE, __, __,
	__, __, __, __, __, __, __, __,
	__, __, __, __, __, __, __, __,
	//              d_quote                              s_quote
	C_SPACE, C_ETC, C_QUOTE, C_ETC, C_DOLAR, C_ETC, C_ETC, C_QUOTE,
	C_LPARB, C_RPARB, C_ASTER, C_PLUS, C_COMMA, C_MINUS, C_POINT, C_SLASH,
	C_ZERO, C_DIGIT, C_DIGIT, C_DIGIT, C_DIGIT, C_DIGIT, C_DIGIT, C_DIGIT,
	C_DIGIT, C_DIGIT, C_COLON, C_ETC, C_ETC, C_ETC, C_ETC, C_QUESM,

	C_AT, C_ABCDF, C_ABCDF, C_ABCDF, C_ABCDF, C_E, C_ABCDF, C_ETC,
	C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC,
	C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC, C_ETC,
	C_ETC, C_ETC, C_ETC, C_LSQRB, C_BACKS, C_RSQRB, C_ETC, C_UNDER,

	C_ETC, C_LOW_A, C_LOW_B, C_LOW_C, C_LOW_D, C_LOW_E, C_LOW_F, C_ETC,
	C_ETC, C_ETC, C_ETC, C_ETC, C_LOW_L, C_ETC, C_LOW_N, C_ETC,
	C_ETC, C_ETC, C_LOW_R, C_LOW_S, C_LOW_T, C_LOW_U, C_ETC, C_ETC,
	C_ETC, C_ETC, C_ETC, C_LCURB, C_ETC, C_RCURB, C_ETC, C_ETC,
}

// The state codes.
const (
	// region JSON

	GO State = iota /* start    */
	OK              /* ok       */
	OB              /* object   */
	KE              /* key      */
	CO              /* colon    */
	VA              /* value    */
	AR              /* array    */
	ST              /* string   */
	ES              /* escape   */
	U1              /* u1       */
	U2              /* u2       */
	U3              /* u3       */
	U4              /* u4       */
	MI              /* minus    */
	ZE              /* zero     */
	IN              /* integer  */
	DT              /* dot      */
	FR              /* fraction */
	E1              /* e        */
	E2              /* ex       */
	E3              /* exp      */
	T1              /* tr       */
	T2              /* tru      */
	T3              /* true     */
	F1              /* fa       */
	F2              /* fal      */
	F3              /* fals     */
	F4              /* false    */
	N1              /* nu       */
	N2              /* nul      */
	N3              /* null     */

	// endregion
	// region JSONPath (root, current, child, rec.des, wildcard, union/slice, script, filter)

	JP /* JSONPath: $ or @     */
	JD /* JSONPath: .          */
	JR /* JSONPath: ..         */
	JW /* JSONPath: *          */
	JO /* JSONPath: [          */
	JH /* JSONPath: [KEY       */
	JA /* JSONPath: [V         */
	JB /* JSONPath: [:         */
	JC /* JSONPath: [:V        */
	JG /* JSONPath: [::        */
	JJ /* JSONPath: [::V       */
	JU /* JSONPath: [,         */
	JL /* JSONPath: ]          */
	JQ /* JSONPath: ?          */
	JS /* JSONPath: (          */
	JF /* JSONPath: )          */
	JK /* JSONPath: dotted key */
	JT /* JSONPath: dotted key quoted */

	// endregion
	// region Script (value, function, function arguments, operation, number, array, object, bool, null, ???)

	SS /* Script: start        */
	SP /* Script: (            */
	SC /* Script: )            */
	SO /* Script: operator     */
	SV /* Script: value        */ // todo: m.b. VA?

	// endregion
)

// The action codes
const (
	cl State = -2 /* colon           */
	cm State = -3 /* comma           */
	qt State = -4 /* quote           */
	bo State = -5 /* bracket open    */
	co State = -6 /* curly br. open  */
	bc State = -7 /* bracket close   */
	cc State = -8 /* curly br. close */
	ec State = -9 /* curly br. empty */
	// jsonpath
	jd State = -10 /* dollar */
	ja State = -11 /* at */
	jq State = -12 /* question mark */
	pc State = -13 /* parentheses close */
	//	todo: function arguments!
	//	todo: implement ErrState and ErrLogic: `a + b * ()` - ErrLogic, `123)` - ErrState. Think about it...
)

// StateTransitionTable is the state transition table takes the current state and the current symbol, and returns either
// a new state or an action. An action is represented as a negative number. A JSONPath text is accepted if at the end of
// the text the state is OK and if the mode is DONE.
// TBD: follow https://github.com/ietf-wg-jsonpath/draft-ietf-jsonpath-base
//             https://datatracker.ietf.org/doc/html/draft-ietf-jsonpath-base
var StateTransitionTable = States{
	/*
	   This is the full state transition table for JSONPath parser.
	                  white                                                    1-9                                                ABCDF                                etc
	            space   |   {   }   [   ]   :   ,   "   \   /   +   -   .   0   |   a   b   c   d   e   f   l   n   r   s   t   u   |   E   (   )   ?   $   @   _   *   |*/
	// region JSON
	/*start  GO*/ {GO, GO, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, JS, __, __, jd, __, __, __, __},
	/*ok     OK*/ {OK, OK, ѢѢ, cc, ѢѢ, bc, ѢѢ, cm, ST, __, __, __, MI, __, ZE, IN, __, __, __, __, __, F1, __, N1, __, __, T1, __, __, __, ѢѢ, __, jq, jd, ja, __, __, __}, // fixme
	/*object OB*/ {OB, OB, ѢѢ, ec, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ST, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*key    KE*/ {KE, KE, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ST, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*colon  CO*/ {CO, CO, ѢѢ, ѢѢ, ѢѢ, ѢѢ, cl, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*value  VA*/ {VA, VA, co, __, bo, __, __, __, ST, __, __, __, MI, __, ZE, IN, __, __, __, __, __, F1, __, N1, __, __, T1, __, __, __, ѢѢ, __, __, JP, JP, ѢѢ, ѢѢ, ѢѢ},
	/*array  AR*/ {AR, AR, co, ѢѢ, bo, bc, ѢѢ, ѢѢ, ST, ѢѢ, ѢѢ, ѢѢ, MI, ѢѢ, ZE, IN, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, F1, ѢѢ, N1, ѢѢ, ѢѢ, T1, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*string ST*/ {ST, ѢѢ, ST, ST, ST, ST, ST, ST, qt, ES, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST, ST},
	/*escape ES*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ST, ST, ST, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ST, ѢѢ, ѢѢ, ѢѢ, ST, ѢѢ, ST, ST, ѢѢ, ST, U1, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*u1     U1*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, U2, U2, U2, U2, U2, U2, U2, U2, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, U2, U2, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*u2     U2*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, U3, U3, U3, U3, U3, U3, U3, U3, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, U3, U3, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*u3     U3*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, U4, U4, U4, U4, U4, U4, U4, U4, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, U4, U4, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*u4     U4*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ST, ST, ST, ST, ST, ST, ST, ST, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ST, ST, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*minus  MI*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ZE, IN, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*zero   ZE*/ {OK, OK, ѢѢ, cc, ѢѢ, bc, ѢѢ, cm, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, DT, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, E1, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, E1, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*int    IN*/ {OK, OK, ѢѢ, cc, ѢѢ, bc, ѢѢ, cm, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, DT, IN, IN, ѢѢ, ѢѢ, ѢѢ, ѢѢ, E1, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, E1, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*dot    DT*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, FR, FR, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*frac   FR*/ {OK, OK, ѢѢ, cc, ѢѢ, bc, ѢѢ, cm, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, FR, FR, ѢѢ, ѢѢ, ѢѢ, ѢѢ, E1, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, E1, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*e      E1*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, E2, E2, ѢѢ, E3, E3, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*ex     E2*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, E3, E3, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*exp    E3*/ {OK, OK, ѢѢ, cc, ѢѢ, bc, ѢѢ, cm, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, E3, E3, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*tr     T1*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, T2, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*tru    T2*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, T3, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*true   T3*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, OK, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*fa     F1*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, F2, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*fal    F2*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, F3, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*fals   F3*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, F4, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*false  F4*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, OK, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*nu     N1*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, N2, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*nul    N2*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, N3, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*null   N3*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, OK, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	// endregion
	// region JSONPath (root, current, child, rec.des, wildcard, union/slice, script, filter)
	/*$|@    JP*/ {__, __, __, cc, JO, bc, __, __, __, __, __, __, __, JD, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, pc, __, __, __, __, __, __},
	/*.      JD*/ {__, __, __, cc, __, bc, __, __, JT, __, __, __, __, JR, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, __, pc, __, __, __, JK, JW, __},
	/*..     JR*/ {__, __, __, cc, __, bc, __, __, JT, __, __, __, __, __, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, JK, __, pc, __, __, __, JK, JW, __},
	/* *     JW*/ {__, __, __, cc, JO, bc, __, __, __, __, __, __, __, JD, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, __, pc, __, __, __, __, __, __},
	/*[      JO*/ {__, __, __, __, __, __, JB, __, JH, __, __, __, ѢѢ, __, ѢѢ, ѢѢ, __, __, __, __, __, __, __, __, __, __, __, __, __, __, JS, __, JQ, __, __, __, __, __},
	/*[KEY   JH*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*[V     JA*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*[:     JB*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*[:V    JC*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*[::    JG*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*[::V   JJ*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*[,     JU*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*]      JL*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*?      JQ*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*(      JS*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*)      JF*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*.key   JK*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*."key" JT*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	// endregion
	// region Script (value, function, function arguments, operation, number, array, object, bool, null, ???)
	/*start  SS*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*(      SP*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*)      SC*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*opera. SO*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	/*value  SV*/ {ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ, ѢѢ},
	// endregion
}

func (s States) GetState(state State, class internal.Class) State {
	return s[state][class]
}

func (s States) GetClass(index byte) Class {
	if index > 128 {
		return C_ETC
	}
	return AsciiClasses[index]
}
