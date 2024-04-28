package newtestcases

import (
	"context"
	"flag"
	"os"
	"strings"
)

var (
	url   = "http://10.2.66.75:28545"
	url_1 = "http://10.2.66.75:18545"

	ctx       = context.Background()
	rootPk    = *rootPrivateKey
	rootPk2   = *root2PrivateKey
	bobPk     = *rootPrivateKey
	builderPk = *builderPrivateKey
	priKey    = os.Getenv("PRIVATE_KEY")

	// setting: root bnb&abc boss
	rootPrivateKey = flag.String("rootpk",
		strings.TrimPrefix(priKey, "0x"),
		"private key of root account")

	root2PrivateKey = flag.String("rootpk2",
		"61bfe9aea17bec5de54a86ad6cb0418f678a2fc8b746cc3901687eaebe1da809",
		"private key of root2 account")

	// root3PrivateKey = flag.String("rootpk3",
	// 	"eb1ee3f15d54f3afcc735ddac56ef8498a006c0bb999a9c267bbf99414698f11",
	// 	"private key of root3 account")

	// root4PrivateKey = flag.String("rootpk4",
	// 	"7540900d280a6df50c6bcaeda216d97df23afb444f82ad840321de853b6bfe9c",
	// 	"private key of root4 account")

	// root5PrivateKey = flag.String("rootpk5",
	// 	"446cdc7ef45999fb635dcbf18acaccd4a796cb7c4fd560b3a6c39b87723e4fc8",
	// 	"private key of root5 account")

	// root6PrivateKey = flag.String("rootpk6",
	// 	"50b9bb6c14ad320ec12b3e21e16296a446059a2453bb9b323a00eb2e051c5eb5",
	// 	"private key of root6 account")

	// root7PrivateKey = flag.String("rootpk7",
	// 	"8fb1b911b16cc94cb2edb8b707c782121c2cf70cd71f2adf2e8bb52bb967a2c4",
	// 	"private key of root7 account")

	builderPrivateKey = flag.String("builderpk",
		"7b94e64fc431b0daa238d6ed8629f3747782b8bc10fb8a41619c5fb2ba55f4e3",
		"private key of builder account")

	validator = flag.String("validator", "0xF474Cf03ccEfF28aBc65C9cbaE594F725c80e12d", "validator address")
)
