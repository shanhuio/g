// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package testkeys

// Pub1 is the public key 1
const Pub1 = "ssh-rsa " +
	"AAAAB3NzaC1yc2EAAAADAQABAAACAQCp0sVYDcRk0LqIC8JMM4a6I8c2cN4Bqw3n" +
	"LsKaPRmzMLR5YaL9UHXUW08s83FXC3+GwrVVgOQ9UCWIs1UI0IDI/TJY2lAK40md" +
	"e0rN+axvz5GXP3ls7OY9HiUaYEwi0UucfjTzZ5GJ7sJmfAsThB0i941iLmClsPsu" +
	"L5mQzcGBWFgwjYmiOfL980CVnZYt/n/NiYqtceE3XCkWgv+d7oWbmAC4WOUcWukc" +
	"GasZSFN3LlbVp6E8B1wmcXjHML+oIPdT3mRWX3cFJLeIKSh0iLOF8u7OqKFou3py" +
	"koT1pT50dAkPnZGTSzGuIjPTOoXo6M1fa8FcdyyN12O52lILF4L+E0Yu1aI/zBSV" +
	"jgtJS60nCDIpJWH0iCWZii4/dSeG0EhK/KxCfXP/03EeU/zh58rOmOl8C2wvP83v" +
	"c2czE321+lu7OxZOy4jUiN4gEPgQfyfcDPSbsXw8V2QzvG2chleCYteoDh/LQbj3" +
	"d5qUGqWz+2JGKmYqZn4wPG8VuQZ3XMf1HtOV044YfEKwcwOOTX4m4YwM2tZ2mM6k" +
	"mOge5XbWclyTCztxqYRCCzESp4xnVDmH1lh8qo8C3nMz8mLxARKBryRzW+tzAVcq" +
	"OaXa7AIbThyOL+GMlywX+cDtntAl5On8DTupaSi0WmUgYgQb5gmQwyeZw6rgQH2A" +
	"RZ00o0OLJQ=="

// Pem1 is the private key 1; its public key is Pub1.
const Pem1 = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
