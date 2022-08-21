package cell

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/xssnick/tonutils-go/address"
)

func TestCell_HashSign(t *testing.T) {
	cc1 := BeginCell().MustStoreUInt(111, 63).EndCell()
	cc2 := BeginCell().MustStoreUInt(772227, 63).MustStoreRef(cc1).EndCell()
	cc3 := BeginCell().MustStoreUInt(333, 63).MustStoreRef(cc2).EndCell()
	cc := BeginCell().MustStoreUInt(777, 63).MustStoreRef(cc3).EndCell()

	b, _ := hex.DecodeString("bb2509fe3cff8f1faae19213774d218c018f9616cd397850c8ad9038db84eaa9")

	if !bytes.Equal(cc.Hash(), b) {
		t.Log(hex.EncodeToString(cc.Hash()))
		t.Log(hex.EncodeToString(b))
		t.Fatal("hash diff")
	}

	pub, priv, _ := ed25519.GenerateKey(nil)
	if !ed25519.Verify(pub, cc.Hash(), cc.Sign(priv)) {
		t.Fatal("sign not match")
	}
}

func TestBOC(t *testing.T) {
	str := "b5ee9c7201021b010003b2000271c000ab558f4db84fd31f61a273535c670c091ffc619b1cdbbe5769a0bf28d3b8fea236865b4312ab35600000625f2d741f0d6773533c74d34001020114ff00f4a413f4bcf2c80b0301510000002629a9a317c878acda0aa0cfacdab9bff8bca840e7d10d8a41d1ee96caf7ac645016af94dfc0160201200405020148060704f8f28308d71820d31fd31fd31f02f823bbf264ed44d0d31fd31fd3fff404d15143baf2a15151baf2a205f901541064f910f2a3f80024a4c8cb1f5240cb1f5230cbff5210f400c9ed54f80f01d30721c0009f6c519320d74a96d307d402fb00e830e021c001e30021c002e30001c0039130e30d03a4c8cb1f12cb1fcbff1213141502e6d001d0d3032171b0925f04e022d749c120925f04e002d31f218210706c7567bd22821064737472bdb0925f05e003fa403020fa4401c8ca07cbffc9d0ed44d0810140d721f404305c810108f40a6fa131b3925f07e005d33fc8258210706c7567ba923830e30d03821064737472ba925f06e30d08090201200a0b007801fa00f40430f8276f2230500aa121bef2e0508210706c7567831eb17080185004cb0526cf1658fa0219f400cb6917cb1f5260cb3f20c98040fb0006008a5004810108f45930ed44d0810140d720c801cf16f400c9ed540172b08e23821064737472831eb17080185005cb055003cf1623fa0213cb6acb1fcb3fc98040fb00925f03e20201200c0d0059bd242b6f6a2684080a06b90fa0218470d4080847a4937d29910ce6903e9ff9837812801b7810148987159f31840201580e0f0011b8c97ed44d0d70b1f8003db29dfb513420405035c87d010c00b23281f2fff274006040423d029be84c6002012010110019adce76a26840206b90eb85ffc00019af1df6a26840106b90eb858fc0006ed207fa00d4d422f90005c8ca0715cbffc9d077748018c8cb05cb0222cf165005fa0214cb6b12ccccc973fb00c84014810108f451f2a7020070810108d718fa00d33fc8542047810108f451f2a782106e6f746570748018c8cb05cb025006cf165004fa0214cb6a12cb1fcb3fc973fb0002006c810108d718fa00d33f305224810108f459f2a782106473747270748018c8cb05cb025005cf165003fa0213cb6acb1f12cb3fc973fb00000af400c9ed5402057fc01817180042bf8e1b0bc5dfcda03e92f9b4b9ffc438595770c0686d91bde674ad610dba9bc66e020148191a0041bf0f895e56f2933fdc5f7c21bc29292fdf0415b7368b9a3eef5bd23ced3021278a0041bf16fc68f92304fb493ca52b5ddefabc42a2131f3e45442b1f2ae45156b2972bea"
	data, _ := hex.DecodeString(str)

	c, err := FromBOC(data)
	if err != nil {
		t.Fatal(err)
	}

	boc := c.ToBOCWithFlags(false)

	if str != hex.EncodeToString(boc) {
		t.Log(str)
		t.Log(hex.EncodeToString(boc))
		t.Fatal("boc not same")
	}
}

func TestBOCWithDecode(t *testing.T) {
	str := "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"
	data, _ := base64.StdEncoding.DecodeString(str)

	c, err := FromBOC(data)
	if err != nil {
		t.Fatal(err)
	}

	boc := c.ToBOCWithFlags(false)

	decodedNew, err := FromBOC(boc)
	if err != nil {
		t.Fatal("boc not parsed")
	}

	if !bytes.Equal(decodedNew.Hash(), c.Hash()) {
		t.Log(str)
		t.Log(hex.EncodeToString(boc))
		t.Fatal("boc not same")
	}
}

func TestSmallBOC(t *testing.T) {
	str := "b5ee9c72010101010002000000"

	c := BeginCell().EndCell()

	boc := c.ToBOCWithFlags(false)

	if str != hex.EncodeToString(boc) {
		t.Log(str)
		t.Log(hex.EncodeToString(boc))
		t.Fatal("boc not same")
	}
}

func TestBOCWithCRC(t *testing.T) {
	str := "b5ee9c7241021b010003b2000271c000ab558f4db84fd31f61a273535c670c091ffc619b1cdbbe5769a0bf28d3b8fea236865b4312ab35600000625f2d741f0d6773533c74d34001020114ff00f4a413f4bcf2c80b0301510000002629a9a317c878acda0aa0cfacdab9bff8bca840e7d10d8a41d1ee96caf7ac645016af94dfc0160201200405020148060704f8f28308d71820d31fd31fd31f02f823bbf264ed44d0d31fd31fd3fff404d15143baf2a15151baf2a205f901541064f910f2a3f80024a4c8cb1f5240cb1f5230cbff5210f400c9ed54f80f01d30721c0009f6c519320d74a96d307d402fb00e830e021c001e30021c002e30001c0039130e30d03a4c8cb1f12cb1fcbff1213141502e6d001d0d3032171b0925f04e022d749c120925f04e002d31f218210706c7567bd22821064737472bdb0925f05e003fa403020fa4401c8ca07cbffc9d0ed44d0810140d721f404305c810108f40a6fa131b3925f07e005d33fc8258210706c7567ba923830e30d03821064737472ba925f06e30d08090201200a0b007801fa00f40430f8276f2230500aa121bef2e0508210706c7567831eb17080185004cb0526cf1658fa0219f400cb6917cb1f5260cb3f20c98040fb0006008a5004810108f45930ed44d0810140d720c801cf16f400c9ed540172b08e23821064737472831eb17080185005cb055003cf1623fa0213cb6acb1fcb3fc98040fb00925f03e20201200c0d0059bd242b6f6a2684080a06b90fa0218470d4080847a4937d29910ce6903e9ff9837812801b7810148987159f31840201580e0f0011b8c97ed44d0d70b1f8003db29dfb513420405035c87d010c00b23281f2fff274006040423d029be84c6002012010110019adce76a26840206b90eb85ffc00019af1df6a26840106b90eb858fc0006ed207fa00d4d422f90005c8ca0715cbffc9d077748018c8cb05cb0222cf165005fa0214cb6b12ccccc973fb00c84014810108f451f2a7020070810108d718fa00d33fc8542047810108f451f2a782106e6f746570748018c8cb05cb025006cf165004fa0214cb6a12cb1fcb3fc973fb0002006c810108d718fa00d33f305224810108f459f2a782106473747270748018c8cb05cb025005cf165003fa0213cb6acb1f12cb3fc973fb00000af400c9ed5402057fc01817180042bf8e1b0bc5dfcda03e92f9b4b9ffc438595770c0686d91bde674ad610dba9bc66e020148191a0041bf0f895e56f2933fdc5f7c21bc29292fdf0415b7368b9a3eef5bd23ced3021278a0041bf16fc68f92304fb493ca52b5ddefabc42a2131f3e45442b1f2ae45156b2972bea32690605"
	data, _ := hex.DecodeString(str)

	c, err := FromBOC(data)
	if err != nil {
		t.Fatal(err)
	}

	boc := c.ToBOC()

	if str != hex.EncodeToString(boc) {
		t.Fatal("boc not same")
	}
}

func TestCell_Hash1(t *testing.T) {
	emptyHash, _ := new(big.Int).SetString("68134197439415885698044414435951397869210496020759160419881882418413283430343", 10)

	if !bytes.Equal(BeginCell().EndCell().Hash(), emptyHash.Bytes()) {
		t.Fatal("empty cell hash incorrect")
		return
	}

	refRef57bitsHash, _ := new(big.Int).SetString("111217512120054409408353636830563617100690868120902564345366655075146083288188", 10)

	if !bytes.Equal(BeginCell().MustStoreUInt(7, 5).MustStoreRef(
		BeginCell().MustStoreRef(
			BeginCell().MustStoreUInt(777777888, 57).EndCell(),
		).EndCell(),
	).EndCell().Hash(), refRef57bitsHash.Bytes()) {
		t.Fatal("refRef57bits cell hash incorrect")
		return
	}
}

func TestCell_InsaneBOC(t *testing.T) {
	// BOC of blocks with index+cache and cell hashes
	str := "b5ee9c72e20201380001000028250000002400cc00ea01c402a603420374039603a503be03d8044804b8050405ac05ec065606a2076e078e0824084208600880089e08bc08da08f80916093009d80a180b0a0b7a0bc70bea0c0e0cba0cda0cfa0d1a0d380d560d740d900dac0dc80e6e0ef20f160f360f820fce0fee100e102e104c106c108c10ac10cc10ec1196121e12841306132413421360137c142014a014ae14bc14ca14d814e614f415021510151e152c153a1548155615a215b015be15cc15da15e815f6160416ba16c816d616e416f21700170e171c172a1738178417a817cc181918c418e419041951199d19bc19da1a271a731a901aae1afb1b471b621b7e1bcb1c171c321cd81d251da81df51e471e921eb21eff1f4b1f6a1f8a1fd71ff620142061208020cd20ec210c2159217821c52211223022da232723ae23fb246024ad24f9257a25c725e42631264e269b26b827052751276c2810285d28dc2929297529c12a0d2ad82b252b442b522b9f2bbc2c092c262c732c922cdf2cfc2d492d662db32dd02e1d2e3a2e872ea42ec22f702fbd3009309e30ac30ba31073114316131ad31ba31c832153222326f327c32c9331533223330337d33c933d633e4343134e634f435aa35b8366c3720372e377b378837d537e237f037fe384b385838a538b238ff390c3959396639b339c03a0d3a1a3a673a743ac13b383bec3c393c463c543ca13cae3cfb3d473d543da13dae3dbc3e093e163ecd3f823fcf3fda3fe0402e405640aa40b740fa410441e84200420e421d422c423c42e043884430443c444844ce458e4614462646ca478b4794481a483648e74988499449a04a264ae64b6c4b7e4c3e4cc44cd64d7b4de94ea84eaf4f354f464fea504b041011ef55aaffffff11000100020003000401a09bc7a98700000000040101485c0d0000000100ffffffff000000000000000062b2ca7f00001a5752b3ec0000001a5752b3ec042d722c2c0004ee1f01485c09014820d8c400000003000000000000002e00050211b8e48dfb4a0eebb004000600071a8a03482793f3b50aaf5c1948a7daea6509532374f902fe6abeff45f620cb8c99cb00130443feeb6ff454dd7b749d28f509060c2cc668963d39733cf0e844a2880964938114c1372a89491186e2103aabb88851625b18674a78929f3c357ca78e98824227016e016e000b000c1489b85d301a5e194c97f1c275d3ead4bd74c6e4f42bc74ff8ef930e594b8dd9887900084a33f6fda87b6952196d2ab3a76db9754f0792826c8eb9650220c062ca9afc456d18d616915f7809f0e2a9af47f9b741c0d1567a8cc87eb371ae985206197a819b2896f2c00109010a010b010c009800001a5752a4a9c401485c0c493505d505412a231566c8853321bb16a8d291b5c091b6b89f7fb0f4470293ca35028007ab3f46b45071071ce8368191ef3b54ec0db1ad04e1ad4b8e254c89e4022581f2cf4db6b621cfac0f967a6e06286bfd400800080008001d43b9aca00250775d8011954fc400080201200009000a0015be000003bcb355ab466ad00015bfffffffbcbd0efda563d0245b9023afe2ffffff1100ffffffff000000000000000001485c0c0000000162b2ca7c00001a5752a4a9c401485c0960000d000e000f0010245b9023afe2ffffff1100ffffffff000000000000000001485c0d0000000162b2ca7f00001a5752b3ec0401485c0960001d001e001f0020284801010190c062d880448c7c066d5e7f424d3c899b5b7618f97ecf7c54e38dd4c8b25a00013213eaca33f4cbadbd5448172264067fbb377adbe16d7c4681fb55312783f206756a318a5c723a464c78c863a17c0802276f4114c4c8edf27dd9eb88622de8445b07016d00118207cb3d36dad8873eb00023008422330000000000000000ffffffffffffffff81f2cf4db6b621cfa828008400222455cc26aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaac23224a27d088a237e001100ac001200ae28480101553f5abf307236fc3eae3b9829943a1b4f895241d07cd2f7fb7106f492a37a4f000222bf000193a937010004ee1f6000034aea52acf0880000d29bdd9fc8200a4106c7d59f0086cde0d05cc99baaddf7a9c5b890aa805e8c96680523524fb4a5a99ff81257e3afd8721409a6a81f22735c0e33400aeb6e7f0a1a5185bb3a3c9deba278be001300142213c3c0000695d4a559e12000b100153201c47dca7c0afc1021a7f8ac40f39fac4e3a34441c180afd0667ef533f08c042258d60f2282e7cde376f0df9fd5ebcc5d8bb19faf6cc57bbbb662485cd76bdce510010000c20004800492211480000d2ba94ab3c2400b30016221162000034aea52acf0900b500172213c480000d2ba94ab3c24000b700182211400000d2ba94ab3c2400b9001922110000034aea52acf09000bb001a22110000034aea52acf09000bd001b2212cc00001a575295678400bf001c2211400000d2ba94ab3c2400c300c4011100000000000000005000213213b82c6f6878b537db3c16dcd5f01397865949762fc38f51cb6aff56594f1062f1edf7f2f7d1549f1d22aab06d9a81c9d49510c3fa6d09dba89fcddb0b8576e4aa016d00118207cb3d37031435feb00068008422330000000000000000ffffffffffffffff81f2cf4dc0c50d7fa828008400223455ee67dd6327a64269ce346e9b1568ef56241f1ba5179cb6d2ac85c5c18646ae12717905e49ef1945ed5b3e03bab93b347378306bc631945e0355ad7507cbada29001b0010cc26aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaac23224a2820ffffb7e010e00ac00ad00ae006bb0400000000000000000a42e0680000d2ba95254e1ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc028480101173ea53338a2f89f176ad2d6129a90977e8c56e325d959bdf92b7171d9070911000423130103e59e9b6d6c439f580024006a00842313010217dee532a6595718002500260084331392afc8f7cde2b08b147ca948f16cc575bbbd4d383188441e172a5cbc617e05249b737c913fe069b3ce1ce7d550cb95f27f4d428faee54cd295c65c5b9ad7e3ae0027000e01014dad03493b489b5800310032008422130100ca31e1e96b10bbc80027006e2213010058b525488168d908006f0028221301003e470bb61928028800290072221100e1caf7c460aa04680073002a221100e1c6b40db0b41708002b0076221100e11e73378043a2080077002c220f00c021469bdc5408002d007a2210680c021412935422007b002e220f00c021383dbc6da8002f007e219dbceaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa81804265fa43649cc3b15c891620f7db60a685fb7054bcc7b8817f05a7e9b73acdede2200b57bfa000034aea549538700302277cff5555555555555555555555555555555555555555555555555555555555555555407aec10e79c0000000000000695d4a92a711804265fa436495d0008000812313010069341ebbd1eac97800330034008422130100e478e48d695dd1e8008500352848010199c07995a55880a2a97e2d1c235a7dfb6b24666e40158c97f7fafb18707d9e01002228480101a40b78c2e39cf3658d455ad29ca1f2bf10cdb87248a8d293e625f4ff821d3ffe001922130100cae245aca13c29a8003600882213010094b93dc1fd640ba800370038221301007e0f065474a37188008b0039221100f6aa376d88c09a280042009f221301007e00c5c8df3ca988008d003a221301007e00c226bc22c348003b0090221301007dfeec4e6f29f348003c0092221301007dfeec4db9634c680093003d221350401f7fbb08ed8733f2003e009621a1bcd9999999999999999999999999999999999999999999999999999999999998200fbfdd8431e351e0e7f7559f2be451b78ede8267e7a8cb24cabd8236aa272459ca146db75357e502000034aea5495385003f227bcff333333333333333333333333333333333333333333333333333333333333333340756c14c3f00000000000000695d4a92a70e00fbfdd8431e351e16d0009800402355ec05b6c5a0ba9cc563f4263b03448b2038580c84db168a6ea499c5ee5d19ddadc8351566240da1dbc38a3b009a009b00412179a062b1fa1362b37a13000080001d81a245901c2c06426d8b4537524ce2f72e8ceed6e41a8ab31206d0ede1c51dc00ebe70af3ef33d8313797763b95f20009d221100f6a49bd9e38f3928004300a1221100ea59d8cb334c9c08004400a3221100ea59cff28ceadfc800a40045220f00c108ba716ce988004600a7219bbd62f8f7bea30f8ab5e9f16c3fb8642b118f56ed1bdc49600dbe5220c8b1af9e040c474f8074f1d14ce894ceb17ddbe1eb5416175db816dcb3d9b86ecacfe305fc3df0baaa80000d2ba95254e1c00047236fcff34517c7bdf5187c55af4f8b61fdc321588c7ab768dee24b006df29106458d7cf21881f480000000000000695d4a92a7110311d3e017f000a900aa00ab220120004a00e2220120005e00c8220120004b00e4220120004c004d220120004e00e822012000f90056220120004f00ea220120005000ec220120005100ee220120005200f0220120005300f2220120005400f4220120005500f628480101a54022f5edbf4beba0648deeb1ffc566c63567b51e34ccd918ab339bdedfec330001220120005700fc220120005800fe22012000ff00592201200101005a220120005b0104220120005c0106220120005d010800b1bd24ee866df51003f6b534d22c17f58175a943d247f628a0d355a187e86e73bd18acb16cc00000000000002a80000001dd3de5878000001dde15bcd058acb28f0000000000000021800000014bff7ab78000001718cb2138a0220120005f00ca22012000cb0060220120006100ce220120006200d0220120006300d222012000d30064220120006500d622012000d70066220120006700da284801015dc67661a1c1ef1294e875961eea55bea7054c7101bc975ed0aa009be971719a000323130103e59e9b818a1aff580069006a00842313010217dee546c430b718006b006c0084284801013c24d22fb5e19c4c445ed87fd2f21aab05a39914a685834b6a1d0cf004891430016b33134d60714c9ba1d20d9c30bf39735d0ad7bfca9eb3ced50edfdae92a4c3339ef4d25e774489adf29b089b2e31b7a68d28886f450600fefbb6f581b3a94423ffc070027000e01014dad035d591ffb5800820083008422130100ca31e1e96b10bbc8006d006e2213010058b525488168d908006f007028480101256e458a2d80eecd799b387aa5e91b156bc583f9f095f0fcea7b2d73cb3d726e00262848010107983f5e2ef514d990c22499abec1ccb4fd222f7489b7febb2044fe24e39318d001a221301003e470bb61928028800710072221100e1caf7c460aa04680073007428480101d405a7172eafa75f0e67e1b0f52d000222f4399c5d7dda059af11c73faf135a300182848010189afd21725efab9748b8b6ef6dcfd892c9f647a939c73831d4f1dbc786b9ec690016221100e1c6b40db0b4170800750076221100e11e73378043a20800770078284801012b525231768abd4f2fc0464b8e5be011db5ebdb3a1ed6d61eebf186e8dc197ba001528480101b81fe7658c95d0f97eabe671d2b147f2da908bd4813119d475886872068353200014220f00c021469bdc54080079007a2210680c021412935422007b007c284801018916fffc4697fa71347a36028c9c0ab1b54bf2c278f9171f7837101f990665f2001128480101a248b81f22333cc28f6b6744e4298aefcd9b6f2dc5d7c99e1da1b28c37f3aa0c0007220f00c021383dbc6da8007d007e219dbceaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa81804265fa43649c9d338c13843e706c68e2455ed70bbf93820f28a394eec0b616e855fce8e6cac000034aea567d807007f284801010143b3d2dd671b2559543155e003f847022e510b3a57afabbca05d4069c327ef000d2277cff5555555555555555555555555555555555555555555555555555555555555555407aec10e79c0000000000000695d4acfb011804265fa436495d0008000812848010164a43970f2007a1da6d6fc81773cc095d1cc270e81359e471f3b03469abeb7b5000c214900000027cbb9d1062954439a83a91f27835fb9d2e3e798910356650c3c493c94623464684000ac28480101afccd0b5d74a6fad14fcb8652b14eea3d3d7ad0226c0961aa9992a1f0c2997c1002322130100e478e4a1873531e80085008628480101a5a7d24057d8643b2527709d986cda3846adcb3eddc32d28ec21f69e17dbaaef0001284801016161938bf6cbbbea618b3c71572f00d9392090c9f15e70d1ef9080503b42229c002322130100cae245c0bf1389a8008700882213010094b93dd61b3b6ba80089008a28480101b4d986be6da4a385ebc4d75e7b6864ddae73e6e6431f711cbb3208dfb35945750024221301007e0f0668927ad188008b008c221100f6aa376d88c09a28009e009f284801012234afc4f9aa54d36f371ab851ada0446bbab533a15161ca6f4afd77d69d89520012221301007e00c5dcfd140988008d008e2848010103f584b35917a807e8a2ec65376c8ed2c6ccecf273ba8cb55dbacfe29a84a6140013221301007e00c23ad9fa2348008f0090221301007dfeec628d0153480091009228480101735d9d54218cf8454bd30b70a53ae0de4219178b78f457b844f50dc11b5719f00013221301007dfeec61d73aac6800930094284801017cc50a2d61f5f6979db4641a9451edade1fb234111081a83333a2b03ce9dc373000a2848010110cac8bd77fa246c2bc95d0269759b612312bbc0438ed5eb7dd53116e55e9020000b221350401f7fbb0df4fd0bf20095009621a1bcd9999999999999999999999999999999999999999999999999999999999998200fbfdd86b59e3de15e7c728be36b6cb9aa8d1663f3d4c1e962ab20459e6e56ae211a233bc1e24f52000034aea567d80500972848010150725eee52e86432f846698a08ac153a67bc9ad9c160130af907c3bef05f29480007227bcff333333333333333333333333333333333333333333333333333333333333333340756c14c3f00000000000000695d4acfb00e00fbfdd86b59e3de16d000980099284801016217f872c99fafcb870f2c11a362f59339be95095f70d00b9cff2f6dcd69d3dd000e2355ec05b6c5a0ba9cc563f4263b03448b2038580c84db168a6ea499c5ee5d19ddadc8351566240da1dbc38a3b009a009b009c28480101fb16d1ca45ecb8d4d1f6b1ac903c630cc06f78334bc9b84bf30585e9422cb887000b284801018f1bd34960aa509ff15ef8c648fdcb942bb7a6c14bda5d4988792ce1c7800bee00062179a062b1fa1362b37a13000080001d81a245901c2c06426d8b4537524ce2f72e8ceed6e41a8ab31206d0ede1c51dc00ebe70af3ef33d831379c7db16df20009d28480101e9988188b13457c31092fcc241e6b801ab7dc39b30a0923557a193630fef257f000a221100f6a49bd9e38f392800a000a128480101c6100af2020b8ed627f48ec736f2cfa52e095b513479105a028f40e152c9586f0016221100ea59d8cb334c9c0800a200a3284801011e20584a9cf50091fb4dddfe4d2e98f8c879438cc843e24314604b44cb6f78580012221100ea59cff28ceadfc800a400a52848010141d3f8101f423d32b2cdb23d98d0f34f83db175e427605500093f4a31cd8df03000b28480101e2bc337ece7f3af5171f3265f44c612fc2fcba87f4b4563dc7fdc3285dd6a44d0008220f00c108ba716ce98800a600a7219bbd62f8f7bea30f8ab5e9f16c3fb8642b118f56ed1bdc49600dbe5220c8b1af9e040c474f800b0df7369fde2be70ee4d2e84a552cab4b035685314de3d2e8e12b8d25aa27ae80000d2ba959f601c000a8284801018e634c5cb159b3914109244a95171c97fe56c7ad67ec709cce94bb067893af680007236fcff34517c7bdf5187c55af4f8b61fdc321588c7ab768dee24b006df29106458d7cf21881f480000000000000695d4acfb0110311d3e017f000a900aa00ab284801017269fb9feb45d719ebdbc3b0816b987bab06f43378dc84dc84d55727905482140002004811fd096c000000000000000000000000000000000000000000000000000000000000000028480101986c49971b96062e1fba4410e27249c8d73b0a9380f7ffd44640167e68b215e8000328480101b4ff459f14a389ff7d6ea967ec8d5329f3cff84a787a7c1fcb6e3d447b6175e5001022bf000193a937010004ee1f6000034aea549538880000d29bdd9fc8200a4106c7d59f0086cde0d05cc99baaddf7a9c5b890aa805e8c96680523524fb4a5a99ff81257e3afd8721409a6a81f22735c0e33400aeb6e7f0a1a5185bb3a3c9deba278be00af00b028480101b20e36a3b36a4cdee601106c642e90718b0a58daf200753dbb3189f956b494b600012213c3c0000695d4a92a712000b100b222012000c500c628480101258d602eaa21d621634dcf86692aeae308ff3cf888f3edafc6a5b21848d732f900182211480000d2ba95254e2400b300b4284801014b01ebcf5425735461aa8b83bae89e70fa21e95d2ee85e57b05dad26c1d6d5300016221162000034aea549538900b500b6284801019523e298bdc5f691343d880493b8a6451f3f941c985a7c4a167ba0e1cdb4599600132213c480000d2ba95254e24000b700b82848010137844b3a6262ee12ef028d6b8968b779c5be75d93a7dc1838aba9408a360ea42000e2211400000d2ba95254e2400b900ba28480101ce05363b2c4d123e6af0a2c3edbe06e05e4b55117180062e69624e00f64ae2a7000c22110000034aea5495389000bb00bc28480101aad2366c8dcad53c429dfbea0cd1c479cc6b989ef981314b1382fd58de7c8afd000b22110000034aea5495389000bd00be284801015af875e56b2c21860165b66c58883a203401027a269c372667df9eb27b476259000a2212cc00001a5752a4a9c400bf00c0284801013eb4f392dec5652b5e530a0922b5533c816263d761c0775349d4265cd85bd761000322110000034aea5495389000c100c222110000034aea52acf09000c300c400a9d00000695d4a92a710000034aea54953880290b818926a0baa0a8254462acd910a6643762d51a5236b81236d713eff61e88e0527946a05000f567e8d68a0e20e39d06d0323de76a9d81b635a09c35a971c4a9913c9284801016e527b5263548e810db4a6e4a1f87c5c20d1c7859c8a5e2113127c954400e7b900012848010192f515d5126f2a2fe83d0204780eb95ac49143dae652a6a618f650997621e09a00013201d970752cf49fb39f7ea882f429ef6a8a5ce3eaf9ff4d35bfba6d93ffc9e2a7368180f9da67bd40c08ccea6759238e3e79b631d22aa1dc395bbcec337be6d4e84000f000c2000e100e222012000c700c822012000c900ca284801012a31c24fa32c257e4912fc641c19d5679a5359fb89020b7fe030e3b104be570f000e22012000cb00cc28480101f8b1119a3146337e09b6a8bbd36f93c78f448ea18e0591af0a3726033c7b5345000c2848010136a19e11f370f7b9cf36d83e0cdb68ea6f02ebe02487e3a5ca264a2baa6cd0d8000c22012000cd00ce22012000cf00d0284801015fb934835d63076694fbd0ca2191f6219a2017bfee8370008d48729374463b93000b22012000d100d228480101d98bd6122c9ac0cd3bb56dd7127ca09ede8c606fcb15de38baa4b47f7fe51f04000922012000d300d4284801011bd2caf8f41ad27d29c9c0d34e01bc0457050b3e5cce1754ec06a5fef2e688af000828480101d3073a3cbcbb4a650ca94a719590f8460e5f51d6081fcea01e7e51cfbc60a444000622012000d500d622012000d700d828480101ab9eb8899afa5c8bb9368d94781cf3ebb2eac7fb017a2b6133f9bb7b0e5d7b85000528480101664d4a2f536f0763f76857ceb94de20b4fad7e80b85675b3d02cf79a319b467e000122012000d900da02012000db00dc284801010f0ee7d20301abe555b4a76a6ada745283de3664ecb09ba794fec4d44db79b68000300b1bd1148769c03366e79fc451b06179e791641db92df408fd46e6ffe7ad24d90a90000000000000000000000000000000000000000000000000000000018a8fe848000000000000014000000006740cdc84000000b659904b9a002012000dd00de00b1bce1814fe876c4b4657ef4585f2f9b9158201e2fc77698c43e68104b88fe4ab2314e7cf70000000000000066800000039c368aef00000040a8c4d343b14e7b67800000000000003b000000040b2e294300000028134825e0c002015800df00e000afbc6d1a11a99e030e7005cac69de1a3cdeea9807743355fbf293f897e485ac868c54bf1b80000000000000124000000091627b10a000000bd5670281ec54bf2a8000000000000011400000006eaf842ee000000a8472576d100afbc6f013e1c5535e8ff36e8381a2acf51990fd66e35d30a40c32f50336512de48c56594fe00000000000001440000000ce66b410a000000d9c852602cc565938800000000000000d40000000de9bcb146000000935ab614cd22012000e300e4284801014356834a429921990e0c978c594d3d7734b7cd96dd66c057aa786ccf4f5a258c000e22012000e500e62848010134a3bc4b6c672bb70328222747df9434117f6ee5953dc9447b945d42862c585f000c22012000e700e822012000f900fa22012000e900ea284801016b54a1df1176e4608f655e63cf51ff3deec889dcf9229cceb100c32a12d3288e000922012000eb00ec284801015d23a19db6638756e655885d66767ab5ca9c7044b8d6263da018629630bda2a5000822012000ed00ee284801014e72d9d9406765832f323a591918a7b8cdb1acdbcbdeab3f2bcbc4c93fbc1f99000922012000ef00f02848010138b4953e5411a45c37899ce1936780cec049051f53d3a16bad52de7c7609476e000622012000f100f228480101a9cf62c6624684d83857a4381d9d525efe2e2d5bafd839d9893f80beeaa43f70000522012000f300f428480101d1a07de28968f21c0c3cfe0e69605731ae68b71969d27dd270f99b4079019904000222012000f500f6284801018d455e04c00fba1289866da39c0793758c8d6a4e41ca2dd5ec2c03c76df76c42000302014800f700f82848010122863015f2ff93c8dfe54f980440d8b3b2bd946770d126f32325fc9c8e37cb2e00020073de48c56594fe000000000290b818000004b380fb4ef8000090e500f4af24c56594fe00000000153ae540000004efd6476f22000098db0ac6ce1f00b0bc914560076a5a3e2f68770608073c7920cd095fd22f6624cdc5631150352c9400000000000000000000000000000000000000000000000000000000629df4ce000000000000004000000004a9d9195a0000002b5c8e819328480101d2fc779057e6ebb9d8413f029999ad6426bdc42b841578ba4f1302565b3bb3b8000a22012000fb00fc22012000fd00fe2848010166827a7a38f195b34597175fc019e4dc80b6f714f0b51a2369270417305747e1000922012000ff010028480101cef4f4863fe525820f12b7c79bbc2276d0fdb7e300c55d70d809bfdef1e801f5000728480101c76ad89ed6929dfcdcbaaa4f5720981c659c02b305514971803760369f7017650006220120010101022848010138e8cca37bb83b168adc87a7037092d90c734af074cace4ae828e0cf82b67b30000422012001030104220120010501062848010150d30f1a3fd6f709f0627934b2b31be081dc25bc52b51a353f850efa09ad224d00022201200107010800b1bd635fefc47229e53988915685933558fd7303875055c5949d871dc2688df796800000000000000000000000000000000000000000000000000000000c54be8f6000000000000007e0000000e610f41dc000000593a932aa9000b1bd24ee866df51003f6b534d22c17f58175a943d247f628a0d355a187e86e73bd18acb16cc00000000000002a80000001dd3de5878000001dde15bcd058acb29fc000000000000021c0000001814ffc698000001752c06e7e202848010160ce52c8bd8ed7f87a7643812f7690467a604fcff9ba1811d065a34a0202d78f000201038020010d00010211011366ea62c6b62574b18990480a15bd04daf2d4d5c8e3413a8f62b0ff533b259b00078201150317cca5687735940043b9aca004010e010f01100247a00d9b55c39995181e04934d61a2baf0f5aa35a4e059bb4c55f309227aee336c95200610011401210103d0400111003fb000000000400000000000000021dcd650010ee6b280087735940043b9aca004010150011301db500cc320a00a42e0680000d2ba95254e000000d2ba95254e0850d8a46888759f13cd17310cd46d6e9e821d494944c72f183259690def72a1a5a15278d7e6ad8df707446b0df93af44f8aca68c8a6c286121567c529e8dd7088800027779c00000000000000000a42e04b159653da0112001343b9aca0021dcd6500200201610114012101064606000125020340400116011702037604011801190297bf955555555555555555555555555555555555555555555555555555555555555502aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaad000000695d4acfb00c1013201340397beb33333333333333333333333333333333333333333333333333333333333333029999999999999999999999999999999999999999999999999999999999999999cf80000695d4acfb00040011a011b011c0397be8517c7bdf5187c55af4f8b61fdc321588c7ab768dee24b006df29106458d7cf029a28be3defa8c3e2ad7a7c5b0fee190ac463d5bb46f71258036f9488322c6be7cf80000695d4acfb0004001270128012901035040011d0103404001210082724a765bd03de7e557d03d203239641e4570953914b82a806fdbc2fa10f8c5e3b7e1cb96b8a440f43d5cc8d5334f00711345b647d43e5837cf367741c70934d23703af7333333333333333333333333333333333333333333333333333333333333333300001a5752b3ec0173fbaacf95f228dbc76f4133f3d46592655ec11b5513922ce50a36dba9abf28100001a5752a4a9c262b2ca7f00014080133011e011f0082724a765bd03de7e557d03d203239641e4570953914b82a806fdbc2fa10f8c5e3b78d942e903175e0ffe9857660057e1e7e904397d1bd95e61a0dc318aaabb1b31f02052030240120013700a0431b9004c4b4000000000000000000960000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003af7333333333333333333333333333333333333333333333333333333333333333300001a5752b3ec0246660377de4bc319a5038b872892ad701fa6a7f2783ecefd07a68cb4fe32f8b700001a5752b3ec0162b2ca7f00014080122012301240101a001250082728d942e903175e0ffe9857660057e1e7e904397d1bd95e61a0dc318aaabb1b31fe1cb96b8a440f43d5cc8d5334f00711345b647d43e5837cf367741c70934d237020f0409283baec018110126013700ab69fe00000000000000000000000000000000000000000000000000000000000000013fccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccd283baec0000000034aea567d800c56594fe40009e42614c107ac00000000000000000640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001035040012a01035040012d008272f46adf2640fb8c30d7d053c9c4084f8daeda7399e261dd404d1578691fcb5bd75f4f3be7c22b3f846e0fc1aa8c5ccb4ed9a7e5acc35369f076448359df9b18eb03af734517c7bdf5187c55af4f8b61fdc321588c7ab768dee24b006df29106458d7cf00001a5752b3ec01e9e3a299d1299d62fbb7c3d6a82c2ebb702db967b370dd959fc60bf87be1755500001a5752a4a9c362b2ca7f00014080133012b012c008272f46adf2640fb8c30d7d053c9c4084f8daeda7399e261dd404d1578691fcb5bd7b43acc176be3014d5645fa082b16fe872be1c3b8428621bc03e49a5bc0cf7db202052030340130013103af734517c7bdf5187c55af4f8b61fdc321588c7ab768dee24b006df29106458d7cf00001a5752b3ec0372f7b0548efa5c07fc58c3051d80c9f650039701e834559eabe556240894c83400001a5752b3ec0162b2ca7f00014080133012e012f008272b43acc176be3014d5645fa082b16fe872be1c3b8428621bc03e49a5bc0cf7db25f4f3be7c22b3f846e0fc1aa8c5ccb4ed9a7e5acc35369f076448359df9b18eb02053030340130013100a042665004c4b400000000000000000030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000069600000009600000004000600000000000519ae84f17b8f8b22026a975ff55f1ab19fde4a768744d2178dfa63bb533e107a409026bc03af7555555555555555555555555555555555555555555555555555555555555555500001a5752b3ec03e61d8ae448b107bedb05342fdb82a5e63dc40bf82d3f4db9d66f6f11005abdfd00001a5752a4a9c362b2ca7f0001408013301340135000120008272010f24a4cdf5d7c0f8497739ff731e5175cd3f10069c838b4445caf821c4cf4e6ca0ffac88c5927e1becceb3949e8aa3388a1cd3e5405413900cce63ed93f19902053030240136013700a041297004c4b40000000000000000002e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005bc00000000000000000000000012d452da449e50b8cf7dd27861f146122afe1b546bb8b70fc8216f0c614139f8e04a5134191"
	data, _ := hex.DecodeString(str)

	c, err := FromBOC(data)
	if err != nil {
		t.Fatal(err)
	}
	_ = c

	// TODO: implement serialize with index
}

func TestCell_Dump(t *testing.T) {
	c := BeginCell().MustStoreInt(-2, 20).EndCell()
	if c.Dump() != "20[FFFFE_]" {
		t.Fatal("wrong dump", c.Dump())
	}

	if c.DumpBits() != "20[11111111111111111110]" {
		t.Fatal("wrong dump", c.DumpBits())
	}
}

func TestVarAddr(t *testing.T) {
	for addrType, str := range map[address.AddrType]string{
		address.NoneAddress: "b5ee9c724101010100030000012094418655",
		address.ExtAddress:  "b5ee9c7241010101000800000b440020406090ae44ae4e",
		address.VarAddress:  "b5ee9c7241010101002900004dd0b000000008012198f3daf0bc973c6958c1e9fc9c65b8ae4e3766a3e89db6e22bae3854ab2b6b5ae33cbe",
		address.StdAddress:  "b5ee9c7241010101002400004389bdc849a2c0204060fdc849a2c0204060fdc849a2c0204060fdc849a2c0204060f01fbc1974",
	} {
		data, _ := hex.DecodeString(str)

		c, err := FromBOC(data)
		if err != nil {
			t.Fatal(err)
		}

		a := c.BeginParse().MustLoadAddr()
		if a.Type() != addrType {
			t.Fatal(addrType, a.Type(), "not correct addr type")
		}

		if hex.EncodeToString(BeginCell().MustStoreAddr(a).EndCell().ToBOC()) != str {
			t.Fatal(addrType, "diff boc")
		}
	}
}
