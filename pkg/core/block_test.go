package core

import (
	"encoding/hex"
	"testing"

	"github.com/CityOfZion/neo-go/pkg/core/transaction"
	"github.com/CityOfZion/neo-go/pkg/crypto"
	"github.com/CityOfZion/neo-go/pkg/io"
	"github.com/stretchr/testify/assert"
)

// Test blocks are blocks from mainnet with their corresponding index.

func TestDecodeBlock1(t *testing.T) {
	data, err := getBlockData(1)
	if err != nil {
		t.Fatal(err)
	}

	b, err := hex.DecodeString(data["raw"].(string))
	if err != nil {
		t.Fatal(err)
	}

	block := &Block{}
	if err := block.DecodeBinary(io.NewBinReaderFromBuf(b)); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, uint32(data["index"].(float64)), block.Index)
	assert.Equal(t, uint32(data["version"].(float64)), block.Version)
	assert.Equal(t, data["hash"].(string), block.Hash().ReverseString())
	assert.Equal(t, data["previousblockhash"].(string), block.PrevHash.ReverseString())
	assert.Equal(t, data["merkleroot"].(string), block.MerkleRoot.ReverseString())
	assert.Equal(t, data["nextconsensus"].(string), crypto.AddressFromUint160(block.NextConsensus))

	script := data["script"].(map[string]interface{})
	assert.Equal(t, script["invocation"].(string), hex.EncodeToString(block.Script.InvocationScript))
	assert.Equal(t, script["verification"].(string), hex.EncodeToString(block.Script.VerificationScript))

	tx := data["tx"].([]interface{})
	minerTX := tx[0].(map[string]interface{})
	assert.Equal(t, len(tx), len(block.Transactions))
	assert.Equal(t, minerTX["type"].(string), block.Transactions[0].Type.String())
	assert.Equal(t, len(minerTX["attributes"].([]interface{})), len(block.Transactions[0].Attributes))
}

func TestTrimmedBlock(t *testing.T) {
	block := getDecodedBlock(t, 1)

	b, err := block.Trim()
	if err != nil {
		t.Fatal(err)
	}

	trimmedBlock, err := NewBlockFromTrimmedBytes(b)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, trimmedBlock.Trimmed)
	assert.Equal(t, block.Version, trimmedBlock.Version)
	assert.Equal(t, block.PrevHash, trimmedBlock.PrevHash)
	assert.Equal(t, block.MerkleRoot, trimmedBlock.MerkleRoot)
	assert.Equal(t, block.Timestamp, trimmedBlock.Timestamp)
	assert.Equal(t, block.Index, trimmedBlock.Index)
	assert.Equal(t, block.ConsensusData, trimmedBlock.ConsensusData)
	assert.Equal(t, block.NextConsensus, trimmedBlock.NextConsensus)

	assert.Equal(t, block.Script, trimmedBlock.Script)
	assert.Equal(t, len(block.Transactions), len(trimmedBlock.Transactions))
	for i := 0; i < len(block.Transactions); i++ {
		assert.Equal(t, block.Transactions[i].Hash(), trimmedBlock.Transactions[i].Hash())
		assert.True(t, trimmedBlock.Transactions[i].Trimmed)
	}
}

func TestHashBlockEqualsHashHeader(t *testing.T) {
	block := newBlock(0)
	assert.Equal(t, block.Hash(), block.Header().Hash())
}

func TestBlockVerify(t *testing.T) {
	block := newBlock(
		0,
		newTX(transaction.MinerType),
		newTX(transaction.IssueType),
	)
	assert.True(t, block.Verify(false))

	block.Transactions = []*transaction.Transaction{
		{Type: transaction.IssueType},
		{Type: transaction.MinerType},
	}
	assert.False(t, block.Verify(false))

	block.Transactions = []*transaction.Transaction{
		{Type: transaction.MinerType},
		{Type: transaction.MinerType},
	}
	assert.False(t, block.Verify(false))
}

func TestBinBlockDecodeEncode(t *testing.T) {
	// transaction taken from mainnet: 2000000
	rawtx := "00000000e5a49e24ee36e972e1bbee16c6897b88050e95e40db157d901cbb68de5243dc93482b51e7ce810eca512afe201768668de5910d4373db067418ad1cf95cd291de424a15a80841e00ec3bd62b5562099d59e75d652b5d3827bf04c165bbe9ef95cca4bf5501fd450140f9ef37e9a31614d0c42aca576d11fcd2ca4cade56143e725ab45e2c7372601e5322a89e0585b44f6f436147be6dc6513ebe781c358abadb1336cadc8f1fdf2e4407fedf529ec4b16ada7fec16efcb377e9c0ea515b12b98a8bed01c385999f8f6121dd5fad32abe4d95dc0c11e9a3a6ce093a7f550b96b779c45f584022bb8a93640d266010bee43509f70c9e7d86cd5037214718de5682abeb42141d1691a1595e5ee188393c26b9ca9f31e4db2d87c3c76869c4b02d081672909268e4d53bcc850401866a84eafd9003c17f1469f1830c5c5f2976da54991f7a1ed292a8af0de2ce202d8f15cb0f362f0ae0ee8bf43886785db45fed0d77b5254503ac105e694a7ac40bfc7166d3495ad4ab540e287ec51afc0569f292e106055b13765d6dacc1ed14807eb63cfeb04b50977c2a64735a4d7496c95f361b773dc58ae29a11b8183f717f1552102486fd15702c4490a26703112a5cc1d0923fd697a33406bd5a1c00e0013b09a7021024c7b7fb6c310fccf1ba33b082519d82964ea93868d676662d4a59ad548df0e7d2102aaec38470f6aad0042c6e877cfd8087d2676b0f516fddd362801b9bd3936399e2103b209fd4f53a7170ea4444e0cb0a6bb6a53c2bd016926989cf85f9b0fba17a70c2103b8d9d5771d8f513aa0869b9cc8d50986403b78c6da36890638c3d46a5adce04a2102ca0e27697b9c248f6f16e085fd0061e26f44da85b58ee835c110caa5ec3ba5542102df48f60e8f3e01c48ff40b9b7f1310d7a8b2a193188befe1c2e3df740e89509357ae120000ec3bd62b0000000080000001e91406db8ce12ff273d9c04cb0f224870d43eeb53c7522a6ad47b0c2ad9614cc0000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc5005ed0b200000000606fe152cafd010c1911b9c7932b77562d13b02f01414054f533865c49ae27ee9898f18c3ad18e2bb27928af9dbfdca44713b666f343e06416fe04869c3b5355934b101ac40caa4910cebdcfc3fec7322d12b0aa4d5bd7232103dfe98cbad29e3116324a5125a32b36250679190f74dca7425c40fff589cc530cac800000012ae528ca6f1c1a740a716603aac1e167cad482d49c41e4a356565cefa50b5dbd0000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc5008d190207000000be87c8c3cb58ca27e5235c5d451474435b895386014140fe20765111b9845d7d842cf09205931ea35e0de5d7078c3831e2b36e41e2298a1cf3d8bc697676ae48a577d70757ecf262085ed314c2d077e0e7a11830dc1b7d232103f2ad7aebe26d46541cb99deb3ff06fd24f0b87efa7e52a1aba6a2d3167be6044acd1015b0600a0724e1809146063795d3b9b3cd55aef026eae992b91063db0db1458f310994f5561f8759ebbdab2e87a05fa4d256553c1087472616e7366657267f91d6b7085db7c5aaf09f19eeec1ca3c0db2c6ecf166aacd14a39b8747f90000000000000000012058f310994f5561f8759ebbdab2e87a05fa4d25650000014140664803b370fb1125fccb44b470a9062d2f1108f4b9c43b7aa35fd56eb52c0564c72e649eefacf82bfb25d252aef5e677307bd0142cecc60e369265b3d053a8d9232103a95af1ee45fe9a2cbad6358e3357da7a13e994361b385c72141ca5a47ef993f6ac800000028190ef017f831707830dd8dc2867572136e7d73bdc13a7421433072e36830c6c00008190ef017f831707830dd8dc2867572136e7d73bdc13a7421433072e36830c6c010001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c600300000000000000494a6c348bdcbf192fb1a88ee41883fcdd11ef37014140a2dc5f161bcb696deeae6f078702627280cfc360d5bf0cae13a417212da2bc87c207d58cca83deb3e6fd76c05bc135367eb520e0138fbfc25e05985ddfa514352321039beeb554e21dfdd414fd07844460d2d1117d26e6b6d18314c535ef25511cee16ac8000000144a1e8637a27d17be06f4d5e4483841c5db27420757350833e970861cf0f6c650500019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500c817a80400000034a3393450b5217977ba2efebad0dd6c37bb0a690141409f6eec2a9983ae782e3d7bdf8ce02f2c9b6db4660ee93197ed8ae85ae9ef424d421ba35a3b70738618b255ab387fe32a0dee1ccd6c0bb8d1e306e3947c383557232103994c84c3cbbc3a75be1045a2ce795d2d5733bb8425f2580e19ff7f179bb8719bac80000001e04c80df4a53b8ba9dca435630f1170d3cc6f15f508e98006ef44bf9f38de82d0000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500bca06501000000d32f28fd47cc176e3d2bf6ac30eb87dfb838ef99014140b8c50990b976e837033eb09da821dae396782bd0df404e1e0675f2db9b7fc6b681b1cb6cb717a546268363f63ae65f89a18769d201b531fee4afd78b2fffb9aa232103ef4c441d2c4cd1a2e3f255366d149553a7f8f799edf1a49a0bc1af0081cc65beac80000001a6f0562becc89fd0ab4c36e2bcc4eedeb6de9fa6bb398be22dafc4e0512b55bc0000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc50001b2c4000000001decb55a65ad7d23a4d474a1e75ff24c9a71dcc4014140175ea6f8e03e7fbf4254ba16ccdd06a6d6c71b7b4d3965bc0913cf02f16ce3180dc278cde426f493225c4e51db28d934da120f5ef2b723585c5aafc8493071e4232103f3215374a4384eb8e2a60635f1e8394ccbe1e21a7b90e0063f5d5ac1836402ceacd100530800fee9e85100000014f9dfa019739920e7ab28544e979c0e6613f54536140a551a5ab4acf6a31dc6fc2a91d05684424a620e53c1087472616e7366657267fb1c540417067c270dee32f21023aa8b9b71abce0001243726662530e8d5576aced9b76cc57b078a9e2857e1ae1bd04869a608657f00000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6001000000000000000a551a5ab4acf6a31dc6fc2a91d05684424a620e014140fd662c8eeaeb3fdc5e9906277c79e31e521b2aa503f0e0684ce876283b27949b2366cd5034a74d7b0b0c47e0648b57a2d569378ee0f82d2bdb3717c4cb3d5c4c2321032204ffdc2173ee04c61f3dbe1ba1ca743106f7cccc7cf6fff86be5dcd6ac63d6ac02000134c42a75a09535438eda8df28e58f8fdc7c5d3329ed6f1d9a845627f7beb53410000000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c60429415000000000046a030db679a5ce51ac083c73b44f0ff6729561c01414065eb2a2a52e6c34ce3a2925ec03b958eeda5cc310468f198ac9c7c111f9a62778214e9b8234c92aef5dba75a29536d6b58b8a074eeb547c53a0e26dc0eaadb86232102670e8d99a3c28d947ddfaa76416a4f53dbcde3ee3d271b1bc78941d5f6d7b829acd1015b0600a0724e1809146063795d3b9b3cd55aef026eae992b91063db0db1458f310994f5561f8759ebbdab2e87a05fa4d256553c1087472616e7366657267f91d6b7085db7c5aaf09f19eeec1ca3c0db2c6ecf1669fb73cc8cde6c72f0000000000000000012058f310994f5561f8759ebbdab2e87a05fa4d25650000014140a1cfd844d2623cfb2b867ec8bf24d3951bd7fcdd5de1e9d59a97b3e240f80020cc7612a3a7dfa32696b7837a1a38f85d27227f76a45b0e6c956abb8b68183177232103a95af1ee45fe9a2cbad6358e3357da7a13e994361b385c72141ca5a47ef993f6ac800000016ab96ecf9f3a32fd8e209d2b5650fb2d12c39bc89493421bdd290d69bbe74d540000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500cf14130400000088eb48d2d1c288c0982ae188b13d5a2a994ab254014140e3269b08dfc3409031a24b7493376a9597eca37a6d6da868fda33005b596accb6e2b7ad4f137d597420d29f080ee9ba545935b47e6358f3f8c1f9118b4e4a085232102e93db927a1f516b37962c4262e78ad56b20e4d349cb1908ee612526ccd0d1118acd10153088005a34c17000000147f808fa952008067570e27dfe0fdc97f647964a51439de68a3e72a015de6ea5543dc896d8609c1c45153c1087472616e7366657267fb1c540417067c270dee32f21023aa8b9b71abce00000000000000000001cd018dd92f5209d47590d74fdd3f982e1ac86d202c12bc948a2890cdb857cff1000002e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c60010000000000000039de68a3e72a015de6ea5543dc896d8609c1c451e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c60a9301b000000000039de68a3e72a015de6ea5543dc896d8609c1c4510141407fad0f5501e9ce3fa9bef9d378257ce233f41d5fd4955e69367ad57392f9feac7608d90553e3e285785cf87fe01b8ab1266d9c73de1e33ceff1c51133735e39b2321028ad2a0b73a8624cc027d80b6ef4af9de5167740aede9cb6401ccc0faef7244bbac0200012d5a5b03a06bbcb891054ae72e64e2d9a3e46077fecf8d2d6af4496771671d910000000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c60a0490100000000008e7778b7ff900a0b82779d0d768cb63b11dcd685014140c54d9d6148503a45de3c037aa6c08ccede953246bdac680a9e9f4ca2414cb02ee130c9a246771191df6bda8796f01f432bfaa464f331137d75865d0c3310f4dc2321036f97d99a1fa298e82ab6ccdedd6febd72fa5b3fa415e09f47354f3cecec51f07acd101530800ba1dd20500000014e31c39cf87027016bf0b925cb6db0699e41d235514bd12389c3eede155f28d6f72c44cb70756aff3f253c1087472616e7366657267fb1c540417067c270dee32f21023aa8b9b71abce0000000000000000000191c12f3ad7a6dc3b8c9cfb6e47d5adb94ef31e76b8dbac538a20a103ba11931f000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c600100000000000000bd12389c3eede155f28d6f72c44cb70756aff3f201414026034a105f0a694f0320b617f112eea7b21a268c54d383e8964e9045e970b8bcc020c3c12e99d560657f2f5874d0a508748bf362bccb0d15f0c9a2c22641ceba2321036200e224b07b5f88b934e36933e5c31826d3b2dc9a5a0a72a555cd5fd70feda5acd1015308803b081117000000141be8765dcb5600ac58a6cc9be39f8d5734c4f1d3148cba604700afc229870ec93ae4f2c25e78f95c4f53c1087472616e7366657267fb1c540417067c270dee32f21023aa8b9b71abce00000000000000000001d4d309fddde1101278ec53443cabb902cbbb2330ce4a16b0876c7f2c411990b5000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6001000000000000008cba604700afc229870ec93ae4f2c25e78f95c4f01414057755879205922936aef9d641f7ac2fd59a9605c01ae5e2cf730ededcc59ef60eaf5d1a5f150c5bd33c3c78b86e441d4d464856bf30ca7b9d89f5299c86d3d7e23210237965fbbff6457d897a9763d6bf97c8f700735e938fa21236aa23e648e641494ac02000109992fadaed0cf0609e42d5fea6fc5057a340893b756c19d73bd0ed92109a0a00000000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c606027000000000000fa1f3b71b6332c5880e1b2808f86e3bf730fbdec0141409ee7a22eee8a502ea3bdd264ba81f36cf7360332456fae92f50d4d89c06cbeda083c61562af8c279e91e490bbabc2e4bee481d85887ff5a36f8dbb26d1308fa22321036980dc27f23a7e3c0ea0fe50f66f6b37a7d21ffaa66d5571ce6d4d317e22a5e7ac80000001cba3499f1640f91f7074593e632d18224e5ca94054c63a72cc3c208fe76a65960100019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500aaa0680b0000007c2f0bc324085be8b5dcf68bc8ed142620b18d850141406020880ded8ae65e6b8af680882f4c3e702cb41f9dcd8a898ea69a32077d05a0503df9d70ff1e71a400dd9a7a2222df9ab348e6cd501473b05add04fb028f960232102f4830c08c254e6e81839396c44c039639f1e30b6e53878269b567d9253590aacac"
	rawtxBytes, _ := hex.DecodeString(rawtx)

	b := Block{}

	r := io.NewBinReaderFromBuf(rawtxBytes)
	err := b.DecodeBinary(r)
	assert.Nil(t, err)
	expected := map[string]bool{ // 18 trans

		"009f61f481f47eb7478e887871e4e744669d461b13d68e04250035260171d706": false,
		"3a62e473c1d67ac561b98e8131f7f7ceded4cd250edb78a6814ec9915930ad93": false,
		"d56a545d2f9400c09d5aa4e8cc37cc994d5a6892f9c30de95ff69a3b647b27a2": false,
		"57f6baa9cb23ce9117d66aee7c31ba6d1e82e140a805db1c5992ada849f6a7c7": false,
		"f700be9c723ed44900ac9d70874b9d8601033bb78883c0a43ab38b5d96c09c11": false,
		"31674ac8553b371ddf06db6a3aef09b8d6f37da03a8cc2868b71044c54ef0034": false,
		"44858de48ec97cea2f823128e9d58981dde11f28a6ebf0a2cb745ea13223dd71": false,
		"317f3ff3768b2aebe3d4866f6e0e8b875cc7937a1b8b5f91be066dc51ed61be2": false,
		"8c24f44f1533567c71e722f49bc7a4d9b323a09e2950fd975291817578119508": false,
		"55a7a738aaee8f7e6d7bcd4c8a38813e57763bff8bfb296418b6cac6d5bfb89a": false,
		"dfa5f84366cf0b48f1b1e9b24a73557e657f6ac21b676528401f5a630aece571": false,
		"5839fbcbbca68aef41dfa9a371222565519626affad6be0977d38a82259480a5": false,
		"6873568cae35e4ce0a7d07ef080ef6eb699b2b9dcbc419fad1c4f645ff8579fc": false,
		"dbb3c0688003bede7e7bc56d2c9d6362b594512ac686820739d963ef91e2eb9e": false,
		"3d12353cb8bae8be928131580e960a82f37ca3ad6957ad22c8cadc1b21b2dd1a": false,
		"8ca87fd5843f000939244151ce027bad5c1f30f1867c7054918b7f9a66b949e8": false,
		"ad088940e45a73e00a3cdb7f3248c67a3f6e5d1f05d4cfd44c4e1f4d26cfef87": false,
		"908a398dd65dfd2aad6c06090c5a71d5e5280746577a6ddd5a1f2c1453f71ead": false,
	}

	hashes := []string{}

	for _, tx := range b.Transactions {
		switch tx.Type {
		case transaction.ContractType:
			hashes = append(hashes, tx.Hash().ReverseString())
		case transaction.MinerType:
			hashes = append(hashes, tx.Hash().ReverseString())
		case transaction.ClaimType:
			hashes = append(hashes, tx.Hash().ReverseString())
		case transaction.InvocationType:
			hashes = append(hashes, tx.Hash().ReverseString())
		}
	}

	assert.Equal(t, len(expected), len(hashes))

	// changes value in map to true, if hash found
	for _, hash := range hashes {
		expected[hash] = true
	}

	// iterate map; all vlaues should be true
	val := true
	for _, v := range expected {
		if v == false {
			val = false
		}
	}
	assert.Equal(t, true, val)

	buf := io.NewBufBinWriter()

	err = b.EncodeBinary(buf.BinWriter)
	assert.Nil(t, err)

	assert.Equal(t, rawtx, hex.EncodeToString(buf.Bytes()))
}

func TestBlockSizeCalculation(t *testing.T) {
	// block taken from mainnet: 0006d3ff96e269f599eb1b5c5a527c218439e498dcc65b63794591bbcdc0516b
	// The Size in golang is given by counting the number of bytes of an object. (len(Bytes))
	// its implementation is different from the corresponding C# and python implementations. But the result should
	// should be the same.In this test we provide more details then necessary because in case of failure we can easily debug the
	// root cause of the size calculation missmatch.

	rawBlock := "00000000ba33df12e8adbf38b6039e79ee91fdb8b1519e2e6154cb59c0653c81769288f4a22492109b7a84077ed7226c28612eb61428ea9ded9bdc952cdfc13deb4172ef85d1115b0bb62300abff35093d19a14a59e75d652b5d3827bf04c165bbe9ef95cca4bf5501fd45014012afae6df64195041e4764b57caa9e27fc2cfc596833163904136ec95816d104b44b3737d0e9f6b1b4445cd3b6a5cc80f6b0935675bc44dba44415eb309832b3404dc95bcf85e4635556a1d618e4ce947b26972992ed74788df5f9501b850ac0b40b7112d1ff30e4ade00369e16f0d13932d1ba76725e7682db072f8e2cd7752b840d12bb7dd45dd3b0e2098db5c67b6de55b7c40164937491fcaca1239b25860251224ead23ab232add78ccccd347239eae50ffc98f50b2a84c60ec5c3d284647a7406fabf6ca241b759af6b71080c0dfad7395632e989226a7e52f8cd2c133aeb2226e6e1aea47666fd81f578405a9f9bbd9d0bc523c3a44d7a5099ddc649feabe5f406188b8ee478731a89beeb76fdbd108eb0071b8f2b8678f40c5a1f387a491314336783255dee8cc5af4bf914dfeaacecc318fc13e02262658e39e8ce0631941b1f1552102486fd15702c4490a26703112a5cc1d0923fd697a33406bd5a1c00e0013b09a7021024c7b7fb6c310fccf1ba33b082519d82964ea93868d676662d4a59ad548df0e7d2102aaec38470f6aad0042c6e877cfd8087d2676b0f516fddd362801b9bd3936399e2103b209fd4f53a7170ea4444e0cb0a6bb6a53c2bd016926989cf85f9b0fba17a70c2103b8d9d5771d8f513aa0869b9cc8d50986403b78c6da36890638c3d46a5adce04a2102ca0e27697b9c248f6f16e085fd0061e26f44da85b58ee835c110caa5ec3ba5542102df48f60e8f3e01c48ff40b9b7f1310d7a8b2a193188befe1c2e3df740e89509357ae140000abff350900000000d101530800e1f5050000000014d8dd86f6d91eb2add2f2d8afeda2184ed94ac33314288017f85d80b889fe02beb5ff203ed9ef538f1653c1087472616e7366657267cf9472821400ceb06ca780c2a937fec5bbec51b900000000000000000220288017f85d80b889fe02beb5ff203ed9ef538f16f0153135323738393433373039313963623362646632650000014140b61b1a8d220c28633fa1a43ef02d334731b16013778664bc28db838d9ab8cfa64f9134fe952ad8a8eb8a6dd9d864055301bcdba4177fd6ee0f52b3f096db2fe4232102bca5c56af0f11e7042f5eaf3d8b2767feb3a8e3ba5668b00e6ea21cb7a215689ac8000000177de54907d16326ff29c3fcd4892afae32043e87ca844929157857093632c4010000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500e9a435000000000849ee84e1f3daaf226589c05de2bf8bcae8c94001414092ba4569663dacd95921318658f8f40662bcff61fdcbbe08da0938a6c93a6d1075b76c86fa5454ca41f762d3b955b8d6755b79ccaf52754169b69e8904f166f8232102e019359f675526fc8505198647e31ed3044ccb0e5cc2ea22fb3bed5420cdf687ac800001ff076e656f2d6f6e6501f3617958f115913192d1de50c8c5ac2c411a4177c7801554d2bf69c27dac86c7010002e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c606400000000000000915fe29d2d3847bce516f31fcd33f0fb1d90573be72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6034ed1a0000000000c85c8dca1fbb7473522109cecaf3acac2e27afc60141406a5ad8c2b6e3783184703d22f3bee39a8b0b6f81477bff61833090361cd52ca5160d66215e6044c31aaed44304a28273c65c9ba736cc75341b45fb18995a6c922321039eae6f12690848807983df6accc1b2929de8582e772be6b3ac084a02a576272eacd10187084d3198000000000008660ceb7bc700000014b3a766ac60afa2990d9251db08138fd1facf07ed0800e3598f01000000209b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc514e420bbe1e5bcdfa1b43e97a57380a6cd8fddfd4c56c1096d616b654f6666657267bd097b2fcf70e1fd30a5c3ef51e662feeafeba0100000000000000000120e420bbe1e5bcdfa1b43e97a57380a6cd8fddfd4c01659c1aa6f2f42f7d0548a9680df197d4dda767434490e8bed3e25dbbac5206cd000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6001000000000000007335f929546270b8f811a0f9427b5712457107e70241400f871fcc8bb58d110ea7b1a4c34390542a925203b147bd68137c9edc79044a657cd5dc00514270d65a77a97c814d1f5090054b6c35f1cde698f765538ccf290c2321023056f0a219758bcd503ff4123a589962003331cb9e14168d649ae7426e3ec26eac4140a4036b311b31e0620fcde6e83ed29e6c7e7fcadf59dcac1ccae352a83608c8117dcd36ca64d5a2ce041b2274481ce65b6d30b463d1ca87a0269177d034ad218d232102a1e6ed9a5cff73ad33b7896465af8e9206eab9c8c75502868b783deb64f232eaacd10179207dcc40f47af5b26c52342d5292eb741b7beaa0d58a9a4a1563441fd0f40c30e1349b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc5b3a766ac60afa2990d9251db08138fd1facf07ed52c10b63616e63656c4f6666657267bd097b2fcf70e1fd30a5c3ef51e662feeafeba01000000000000000001202e5596a5c23eff8907aa180201c5a6f53c041dd101a9822ade57f6f9b6e999be06961af808137f3c9051b6c674bcc9e8aec8c1b6f1000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6001000000000000007335f929546270b8f811a0f9427b5712457107e7024140b5b58794f8c4f2493a2cff25aa7546a79b6f9e1f279911b2d15a07926084705aeaf64d18454f40feda8ddb9096745eb2e16a208320e9122007fac350892cb3d823210308ec2156f3366339c54c59e4c0342888665abdd76ccbd6b2020961225ccfa3f4ac41404f947ca69b7a2ec0926989f30f3d3d986d488b7c01270eee56326b2f0be856353448062a5d2d73ce34bfdacec3af42b082f1c5186fc063b38c321c1a1eb37157232102a1e6ed9a5cff73ad33b7896465af8e9206eab9c8c75502868b783deb64f232eaacd101792084e25c75cef1e92d39333408dae6e31799d6316b2d908aa094d0dc18f137484c349b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc5b3a766ac60afa2990d9251db08138fd1facf07ed52c10b63616e63656c4f6666657267bd097b2fcf70e1fd30a5c3ef51e662feeafeba01000000000000000001202e5596a5c23eff8907aa180201c5a6f53c041dd101ead2059c2cf4101cb02a1dae5875dc993189d047a84f05c080318779453b3c0e000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6001000000000000007335f929546270b8f811a0f9427b5712457107e7024140c4c924181d2b1516467f20fd3cf4b75ce5dd75d916318b77cd9ebd4bf9a3fb3bef47151adb01e0828bc127db6859f11890f86bc5bae2980cba5fcefba0d7edf523210308ec2156f3366339c54c59e4c0342888665abdd76ccbd6b2020961225ccfa3f4ac41407beaabc59bbc4db241aca812b624befcc595d3f8299e996a7d660decd3cb162afc6d810d5e3001aae5874c42b7031b1445f8562c209a2e2726ca1316dffd8867232102a1e6ed9a5cff73ad33b7896465af8e9206eab9c8c75502868b783deb64f232eaacd101530800e1f5050000000014d8dd86f6d91eb2add2f2d8afeda2184ed94ac33314288017f85d80b889fe02beb5ff203ed9ef538f1653c1087472616e7366657267f91d6b7085db7c5aaf09f19eeec1ca3c0db2c6ec00000000000000000220288017f85d80b889fe02beb5ff203ed9ef538f16f015313532373839343337373637356464383030353031000001414061609a0460a3ccfe1a9cb5db9f75811e08d52328f291a1b848ec607718be0a37206a90e3a81908c4b71ec859b684e493c088e640b2e2d471bd370aae50bdf160232102bca5c56af0f11e7042f5eaf3d8b2767feb3a8e3ba5668b00e6ea21cb7a215689ac80000001128218ba3c40a03066c862b0eaad5c06a39a1ae653e08f69e70f528ec4e18dc60500019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500a3e11100000000db2395d79fbd27d7b93b62cccd0fe0afc15ff80b0141404d953b03a1f53911fb1524a9f12126790e5e8468c05ed24ac3f71b881623510fcbe2f5497a4e8d1f1d1748774b9b3487be5f45dbd91bfb697a27ae2d2f2d9bb32321030ab39b99d8675cd9bd90aaec37cba964297cc817078d33e508ab11f1d245c068acd101792017a323cf2b0c650243be29d1b2c3ca69c141fdd7d72a36242c225f175c818a46349b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc5b3a766ac60afa2990d9251db08138fd1facf07ed52c10b63616e63656c4f6666657267bd097b2fcf70e1fd30a5c3ef51e662feeafeba01000000000000000001202e5596a5c23eff8907aa180201c5a6f53c041dd1010a190ead75c8235b2db69d9eb00040209ce719c3450d90bca893bcfcdbe5e58c000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6001000000000000007335f929546270b8f811a0f9427b5712457107e7024140b1dbd06a250ee76dc773e7efc915caea48732ca10f813d1027073dbba57032752c0ac7ed30fdb8b18df0abaf5831bd86e444785445253099b46c271676c4a16423210308ec2156f3366339c54c59e4c0342888665abdd76ccbd6b2020961225ccfa3f4ac4140a3b256ac1dd57de358bc1a45b632b3dd41f07d4c6d6e9ffe33fd0cf7c44e29673100bc8a8760727b46bf96642d6dd4a9c4edf657fb2a27b83f564a5b91408ad9232102a1e6ed9a5cff73ad33b7896465af8e9206eab9c8c75502868b783deb64f232eaacd1014e03a0860114d023de91710f63a0259e8d95d1f6563e1572783c14b18f53a903d7873b1453a748c3f80787eca2e30f53c1087472616e7366657267187fc13bec8ff0906c079e7f4cc8276709472913000000000000000004fee36e656f2d6f6e652d696e766f6b653a7b22636f6e7472616374223a22307831333239343730393637323763383463376639653037366339306630386665633362633137663138222c226d6574686f64223a227472616e73666572222c22706172616d73223a5b5b2266726f6d222c22307830666533613265633837303766386333343861373533313433623837643730336139353338666231225d2c5b22746f222c22307833633738373231353365353666366431393538643965323561303633306637313931646532336430225d2c5b2276616c7565222c22302e303031225d5d7dff076e656f2d6f6e65ff0a3233313635343031383420b18f53a903d7873b1453a748c3f80787eca2e30f0000014140910c4b3be37fe09052215f99ef546342c21ded22ef862c165ff0da3ff1087274d0fafb5bbac25aa927dc477b239df732d95afe782a46caeaca669ae096659f13232102c2dbc83931d5e550b95ceab8a94c6af37735fe2aa4e9fb217bce46001937b2f1acd1014e03a0860114de33c5d07f933c0f90da952cdfe2677d3f9a24d714f26e3a75379dc20c8642f3f11e20af76b12065f453c1087472616e7366657267aa67d0447c61bdddc4d1690d2269d772f9c37795000000000000000004fee36e656f2d6f6e652d696e766f6b653a7b22636f6e7472616374223a22307839353737633366393732643736393232306436396431633464646264363137633434643036376161222c226d6574686f64223a227472616e73666572222c22706172616d73223a5b5b2266726f6d222c22307866343635323062313736616632303165663166333432383630636332396433373735336136656632225d2c5b22746f222c22307864373234396133663764363765326466326339356461393030663363393337666430633533336465225d2c5b2276616c7565222c22302e303031225d5d7dff076e656f2d6f6e65ff0933363133393630393120f26e3a75379dc20c8642f3f11e20af76b12065f40000014140b1630c1546547df4b69fdf6c68fcbf9a3f1de2787d0ea7237aea639ab0a227a092147e47194595c0d960df0d209773da32808da767c15730926bcf844ed12f702321021a9f5ed87fe58e7a366c20975d4112698e4c0ccb3ba9cbce0400a482ecf99b67acd1014e03a08601148ba89ee5ab5c4975e0e12f88a8ce4aa8928108c9145e5c739e5d4b5a29af596c2cc162f1facded25fe53c1087472616e7366657267952d12a025325e56a4cb3ba2d469b1e23c7c77a0000000000000000004fee36e656f2d6f6e652d696e766f6b653a7b22636f6e7472616374223a22307861303737376333636532623136396434613233626362613435363565333232356130313232643935222c226d6574686f64223a227472616e73666572222c22706172616d73223a5b5b2266726f6d222c22307866653235656463646661663136326331326336633539616632393561346235643965373335633565225d2c5b22746f222c22307863393038383139326138346163656138383832666531653037353439356361626535396561383862225d2c5b2276616c7565222c22302e303031225d5d7dff076e656f2d6f6e65ff0a33383237333138333238205e5c739e5d4b5a29af596c2cc162f1facded25fe000001414050f556b29417eb44b66b283a72ff33b697f49b2b6d1520eb4fe5a9974632c1b1586b3e7f0aafac4707d65f89daa86592c8ad98c78016769a5de904330d02d2ed23210201008fe0ffcdab73b598c89c6ae2b46d90de38287abd7dd50a325d0bfb2469d5acd1014e03a086011444d65fda3f2062502c03c2bfd85c700a0d046fae149d12fcef2c830d73a570e9b89857962fcf3a619a53c1087472616e7366657267187fc13bec8ff0906c079e7f4cc8276709472913000000000000000004fee36e656f2d6f6e652d696e766f6b653a7b22636f6e7472616374223a22307831333239343730393637323763383463376639653037366339306630386665633362633137663138222c226d6574686f64223a227472616e73666572222c22706172616d73223a5b5b2266726f6d222c22307839613631336163663266393635373938623865393730613537333064383332636566666331323964225d2c5b22746f222c22307861653666303430643061373035636438626663323033326335303632323033666461356664363434225d2c5b2276616c7565222c22302e303031225d5d7dff076e656f2d6f6e65ff0a31353931373638373437209d12fcef2c830d73a570e9b89857962fcf3a619a00000141409ec506aab7045da733d02c8c9ebaa615d9aa02fa8c4d8eb35ad08a3b4843167cdd059094f492f77e68efa2b66d87305e13b3f08a326cf4b3dadaa29aaa64b30a23210275699ef1532219ee8e9d2d8a75b6b96b7cb4d5f9788391de31daceebed9ac8edacd1014e03a086011444bb85601b2d8c5247a1999dfb18ee7928e10cdf1479ed35e989051c8f2404b869c17ce75912de78b953c1087472616e7366657267aa67d0447c61bdddc4d1690d2269d772f9c37795000000000000000004fee36e656f2d6f6e652d696e766f6b653a7b22636f6e7472616374223a22307839353737633366393732643736393232306436396431633464646264363137633434643036376161222c226d6574686f64223a227472616e73666572222c22706172616d73223a5b5b2266726f6d222c22307862393738646531323539653737636331363962383034323438663163303538396539333565643739225d2c5b22746f222c22307864663063653132383739656531386662396439396131343735323863326431623630383562623434225d2c5b2276616c7565222c22302e303031225d5d7dff076e656f2d6f6e65ff0a333738313933373030382079ed35e989051c8f2404b869c17ce75912de78b9000001414081a2d21af431df3345478836ae24dc7714ce647daa1daac3a99ed17645a8c86a981d86c5d2280af45259a1cf4a85c115ecbe8ff5a398c0bea23c218c0c8c26f02321020be10b4bddffe752a7cfaa16d1718ce6da460608ca41ad9bf7ff66dc5f60c860acd10179201322db68df23cee8da600e7cf5875a8d919bf8e318bbe7afbcd1fdbd5758d67b349b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc5b3a766ac60afa2990d9251db08138fd1facf07ed52c10b63616e63656c4f6666657267bd097b2fcf70e1fd30a5c3ef51e662feeafeba01000000000000000001202e5596a5c23eff8907aa180201c5a6f53c041dd10181c3670591e333a19a14aae556d211f3490163380c1f28ef55b0bb0e4a8aec1a000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6001000000000000007335f929546270b8f811a0f9427b5712457107e7024140dd2aa6075329f43375655c83b04d343974fb13908888fa6b242b599c789c1cf635b530b6834093332720f93d185b450884951acfb261aec61fe9057efbdeb4b523210308ec2156f3366339c54c59e4c0342888665abdd76ccbd6b2020961225ccfa3f4ac4140f152dc9454ebf518bc0ab8192543add5844ad7007177986a3dd038db313740ea72abb5d49866caa4d470d75182a7e715fd2b9a6e345ef2568bf82e92878743db232102a1e6ed9a5cff73ad33b7896465af8e9206eab9c8c75502868b783deb64f232eaac8000000198f4ec54f46fd86c0955accb7bda7ff06ffc10c45651206fbcbfeabe78c597750000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500e40b54020000002fb7e583c973498fef06f317a9762507e7ac0306014140c4a7388365a4e8a93e2229fd617805902773a8ef7354e95487819eec59379a6a8fb045ed126025db151336e2393eaf75dd524f1ebe872f92185f8d1a6a71a6a5232102b37eaec8631a5bb8579d4ba268e2fbc9c81b555f08558bba6eccb9d1448332c5ac02000155810e77f1dd89622915f3fceb051ac8abeaa2e7ea7a2c38e0744137c3c8982b0000000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c605e1502000000000024fdf844d637cf5a60a9744f1ab8ab8d6f2c18cf0141407dafb244fceee305d477892a21ca8cbd31906fbcc233e0bfccd0fe91d4f6dafdd067f7301cd3e8c7c6197f71615df4b7cedb987eb3ebff79ad729bb8768e1e50232102f8a4a73deefad4c114039b854cc4e974c409f30720e2567d15dea2204b8c45f7acd1012000c108776974686472617769bd097b2fcf70e1fd30a5c3ef51e662feeafeba01000000000000000003a15100000000000000000000000000000000000000000000000000000000000000a232e125258b7db0a0dffde5bd03b2b859253538ab000000000000000000000000a48098c835c493eb3b1967e44150630cc6435e564e00000000000000000000000001dfeee6c60db6403b5538e86120f01e60dcea749c73ce36b9058f050be8b7d2c9000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c600100000000000000bd097b2fcf70e1fd30a5c3ef51e662feeafeba010102000000d101530800e1f5050000000014d8dd86f6d91eb2add2f2d8afeda2184ed94ac33314288017f85d80b889fe02beb5ff203ed9ef538f1653c1087472616e7366657267f56c89be8bfcdec617e2402b5c3fd5b6d71b820d00000000000000000220288017f85d80b889fe02beb5ff203ed9ef538f16f0153135323738393433383433393264623331653662630000014140654af3a0cc69faa3dd42ff76c4012aa9c72e269bba004d6e910f195e33b2ecae980be4531a3677f27d3c90f4196632790997078bca4f8471c6db43b55928c3ef232102bca5c56af0f11e7042f5eaf3d8b2767feb3a8e3ba5668b00e6ea21cb7a215689ac"
	rawBlockBytes, _ := hex.DecodeString(rawBlock)

	b := Block{}

	r := io.NewBinReaderFromBuf(rawBlockBytes)
	err := b.DecodeBinary(r)
	assert.Nil(t, err)

	expected := []struct {
		ID            string
		Type          string
		Size          int
		Version       int
		InputsLen     int
		OutputsLen    int
		AttributesLen int
		WitnessesLen  int
	}{ // 20 trans
		{ID: "f59b04d8e6526684b94b5f8cdbdf691feaff5d45e9aa8e2325a668f1b9130786", Type: "MinerTransaction", Size: 10, Version: 0, InputsLen: 0, OutputsLen: 0, AttributesLen: 0, WitnessesLen: 0},
		{ID: "7463345f771e70019185d72fa5bd00fbb4f26735daae398ecc6540419332d81e", Type: "InvocationTransaction", Size: 244, Version: 1, InputsLen: 0, OutputsLen: 0, AttributesLen: 2, WitnessesLen: 1},
		{ID: "cf3aeda21d320ec9b49d322f2b88fea21aa7e9bf243c1e02dfe08f5cc82b74b0", Type: "ContractTransaction", Size: 202, Version: 0, InputsLen: 1, OutputsLen: 1, AttributesLen: 0, WitnessesLen: 1},
		{ID: "07e502d13ae6255cfabbc9ee2f78a48fc1c43a4f7f713f128342db721bc01af5", Type: "ContractTransaction", Size: 271, Version: 0, InputsLen: 1, OutputsLen: 2, AttributesLen: 1, WitnessesLen: 1},
		{ID: "0de0baf53136c188bdd179fed9530dfb7dd80697fd59e47ffe294db4f421eb67", Type: "InvocationTransaction", Size: 469, Version: 1, InputsLen: 1, OutputsLen: 1, AttributesLen: 1, WitnessesLen: 2},
		{ID: "233c8b00ab6a43aafae7fcc2be47fc46493185bb3376160b5809cb745aee3329", Type: "InvocationTransaction", Size: 455, Version: 1, InputsLen: 1, OutputsLen: 1, AttributesLen: 1, WitnessesLen: 2},
		{ID: "8bf3ae0c692fc830753029fcb6575625ea8181b444cffcbe38404a28b77b3856", Type: "InvocationTransaction", Size: 455, Version: 1, InputsLen: 1, OutputsLen: 1, AttributesLen: 1, WitnessesLen: 2},
		{ID: "c5285e1460191e1ca7fc07e3c26c5facebb033d56b63b7a41ebf11f2a1cb4306", Type: "InvocationTransaction", Size: 244, Version: 1, InputsLen: 0, OutputsLen: 0, AttributesLen: 2, WitnessesLen: 1},
		{ID: "dce8b5f6dc093a910e405a230e9b7d546688d411cf960c8a1cc7d386d89b56d6", Type: "ContractTransaction", Size: 202, Version: 0, InputsLen: 1, OutputsLen: 1, AttributesLen: 0, WitnessesLen: 1},
		{ID: "4cce087cadfa99c2adeaaf1916ada025db124cef8f05d4535b0ad8047ef7d29e", Type: "InvocationTransaction", Size: 455, Version: 1, InputsLen: 1, OutputsLen: 1, AttributesLen: 1, WitnessesLen: 2},
		{ID: "67df57a20c9d3b2942925f2c66fdc15a21be2c229a22122f6acbdac4dd10bf0a", Type: "InvocationTransaction", Size: 466, Version: 1, InputsLen: 0, OutputsLen: 0, AttributesLen: 4, WitnessesLen: 1},
		{ID: "ad51030b30e016293caed92781b3bb3f993f86c15ab1153582f658d603fe23db", Type: "InvocationTransaction", Size: 465, Version: 1, InputsLen: 0, OutputsLen: 0, AttributesLen: 4, WitnessesLen: 1},
		{ID: "1db2d62ad3530f1ae6ca7bd95e766beaff97058681f0e203d8744d7bba065012", Type: "InvocationTransaction", Size: 466, Version: 1, InputsLen: 0, OutputsLen: 0, AttributesLen: 4, WitnessesLen: 1},
		{ID: "7d7bb6f0db6a71aca85fc9267fa6a59654b00b5f778a39c27214c68f11950f61", Type: "InvocationTransaction", Size: 466, Version: 1, InputsLen: 0, OutputsLen: 0, AttributesLen: 4, WitnessesLen: 1},
		{ID: "1d534dcf1ce63a9ea9328eab309891ea2a0a5cb11e95cabf22860ee1fb649521", Type: "InvocationTransaction", Size: 466, Version: 1, InputsLen: 0, OutputsLen: 0, AttributesLen: 4, WitnessesLen: 1},
		{ID: "8f4c9089871a4ad0076b27c061395079e0862f685f27e4bc01b7bac67b0cf8d0", Type: "InvocationTransaction", Size: 455, Version: 1, InputsLen: 1, OutputsLen: 1, AttributesLen: 1, WitnessesLen: 2},
		{ID: "52b4653fca02e8042092456490036f0b9b18b339f65d9c334e7e9d2b4599f8db", Type: "ContractTransaction", Size: 202, Version: 0, InputsLen: 1, OutputsLen: 1, AttributesLen: 0, WitnessesLen: 1},
		{ID: "2b4854b1f46c9af0eb06587fd375355adfeea3c8f5921295421251af93d737e1", Type: "ClaimTransaction", Size: 203, Version: 0, InputsLen: 0, OutputsLen: 1, AttributesLen: 0, WitnessesLen: 1},
		{ID: "61e95e5b14625e897423670dfc3babf021d6b99ca2a73203dd1ac2604a2daadf", Type: "InvocationTransaction", Size: 244, Version: 1, InputsLen: 1, OutputsLen: 1, AttributesLen: 3, WitnessesLen: 1},
		{ID: "b361dfec8c2cde980b340d2c3ec63cecaea634f91b6d76f24a586aa60fbde483", Type: "InvocationTransaction", Size: 244, Version: 1, InputsLen: 0, OutputsLen: 0, AttributesLen: 2, WitnessesLen: 1},
	}

	for i, tx := range b.Transactions {
		txID := tx.Hash()
		assert.Equal(t, expected[i].ID, txID.ReverseString())

		assert.Equal(t, expected[i].Size, tx.Size())
		assert.Equal(t, expected[i].Type, tx.Type.String())
		assert.Equal(t, expected[i].Version, int(tx.Version))
		assert.Equal(t, expected[i].InputsLen, len(tx.Inputs))
		assert.Equal(t, expected[i].OutputsLen, len(tx.Outputs))
		assert.Equal(t, expected[i].AttributesLen, len(tx.Attributes))
		assert.Equal(t, expected[i].WitnessesLen, len(tx.Scripts))
	}

	assert.Equal(t, len(expected), len(b.Transactions))

	// Block specific tests
	assert.Equal(t, 0, int(b.Version))
	assert.Equal(t, "f4889276813c65c059cb54612e9e51b1b8fd91ee799e03b638bfade812df33ba", b.PrevHash.ReverseString())
	assert.Equal(t, "ef7241eb3dc1df2c95dc9bed9dea2814b62e61286c22d77e07847a9b109224a2", b.MerkleRoot.ReverseString())
	assert.Equal(t, 1527894405, int(b.Timestamp))
	assert.Equal(t, 2340363, int(b.Index))

	nextConsensus := crypto.AddressFromUint160(b.NextConsensus)
	assert.Equal(t, "APyEx5f4Zm4oCHwFWiSTaph1fPBxZacYVR", nextConsensus)

	assert.Equal(t, "4012afae6df64195041e4764b57caa9e27fc2cfc596833163904136ec95816d104b44b3737d0e9f6b1b4445cd3b6a5cc80f6b0935675bc44dba44415eb309832b3404dc95bcf85e4635556a1d618e4ce947b26972992ed74788df5f9501b850ac0b40b7112d1ff30e4ade00369e16f0d13932d1ba76725e7682db072f8e2cd7752b840d12bb7dd45dd3b0e2098db5c67b6de55b7c40164937491fcaca1239b25860251224ead23ab232add78ccccd347239eae50ffc98f50b2a84c60ec5c3d284647a7406fabf6ca241b759af6b71080c0dfad7395632e989226a7e52f8cd2c133aeb2226e6e1aea47666fd81f578405a9f9bbd9d0bc523c3a44d7a5099ddc649feabe5f406188b8ee478731a89beeb76fdbd108eb0071b8f2b8678f40c5a1f387a491314336783255dee8cc5af4bf914dfeaacecc318fc13e02262658e39e8ce0631941b1", hex.EncodeToString(b.Script.InvocationScript))
	assert.Equal(t, "552102486fd15702c4490a26703112a5cc1d0923fd697a33406bd5a1c00e0013b09a7021024c7b7fb6c310fccf1ba33b082519d82964ea93868d676662d4a59ad548df0e7d2102aaec38470f6aad0042c6e877cfd8087d2676b0f516fddd362801b9bd3936399e2103b209fd4f53a7170ea4444e0cb0a6bb6a53c2bd016926989cf85f9b0fba17a70c2103b8d9d5771d8f513aa0869b9cc8d50986403b78c6da36890638c3d46a5adce04a2102ca0e27697b9c248f6f16e085fd0061e26f44da85b58ee835c110caa5ec3ba5542102df48f60e8f3e01c48ff40b9b7f1310d7a8b2a193188befe1c2e3df740e89509357ae", hex.EncodeToString(b.Script.VerificationScript))
	assert.Equal(t, "0006d3ff96e269f599eb1b5c5a527c218439e498dcc65b63794591bbcdc0516b", b.Hash().ReverseString())

	buf := io.NewBufBinWriter()

	err = b.EncodeBinary(buf.BinWriter)
	assert.Nil(t, err)
	benc := buf.Bytes()
	// test size of the block
	assert.Equal(t, 7360, len(benc))
	assert.Equal(t, rawBlock, hex.EncodeToString(benc))
}
