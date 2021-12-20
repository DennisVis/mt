// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package mt

// MT940 represents a Customer Statement Message.
// It's based on the spec here: https://www2.swift.com/knowledgecentre/publications/us9m_20210723/1.0?topic=mt940.htm
type MT940 struct {
	Base
	Reference                     string          `mt:"20,M,16x"`
	AccountIdentification         string          `mt:"25,M,2!c26!n|8!c/12!n"`
	StatementNumberSequenceNumber string          `mt:"28C,M,5!n(/3!n)"`
	OpeningBalance                Balance         `mt:"60F,M,dive"`
	StatementLines                []StatementLine `mt:"61,O,dive"`
	AccountOwnerInformation       []string        `mt:"86,O,6*65x"`
}
