# Changelog

## [0.9.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2-v0.8.0...tabelogo-v2-v0.9.0) (2025-12-24)


### Features

* **auth:** implement full session management (Phase 1.3) ([0c015e9](https://github.com/Leon180/tabelogo-v2/commit/0c015e9b03998d606765e03f8e9e491d555d2031))
* **auth:** implement full session management (Phase 1.3) ([82f7142](https://github.com/Leon180/tabelogo-v2/commit/82f7142e54d7a5fb057ac05de6c98519de516161))
* **config:** migrate JWT configuration to environment variables (Phase 2.5) ([dba0650](https://github.com/Leon180/tabelogo-v2/commit/dba0650236bc00dbdb6c0fe4f11db1d594fb1dcb))
* **docs:** add integration tests and deployment documentation (Phase 2.6) ([d64a955](https://github.com/Leon180/tabelogo-v2/commit/d64a9559d9f001798903303fcad214d65f24c805))
* **frontend:** add authentication to Spider Service API client ([52435c3](https://github.com/Leon180/tabelogo-v2/commit/52435c37dc2317a1104e989fcb024d0375de4d84))
* **frontend:** add register page and auth integration docs ([e8075da](https://github.com/Leon180/tabelogo-v2/commit/e8075da9c1398c5ce46e0538353ea21c75cc4214))
* **map:** integrate auth middleware with Optional strategy (Phase 2.4) ([7dd3eff](https://github.com/Leon180/tabelogo-v2/commit/7dd3effe805253c20c8dbdfac6cd831cd3d11e98))
* **middleware:** update auth middleware with session validation (Phase 2.1) ([2c32769](https://github.com/Leon180/tabelogo-v2/commit/2c3276967393d18694aa04fbda97533e4b2e91ce))
* **prometheus:** add spider-service to scrape configuration ([a0580c5](https://github.com/Leon180/tabelogo-v2/commit/a0580c5ae4e09a43edff556764d60f7e4ad9f667))
* **restaurant:** integrate auth middleware with mixed public/protected routes (Phase 2.3) ([7d40d50](https://github.com/Leon180/tabelogo-v2/commit/7d40d50160a59c5e62a63bb8648f382e0543209a))
* **spider:** integrate auth middleware (Phase 2.2) ([e275bf4](https://github.com/Leon180/tabelogo-v2/commit/e275bf46443b0482f52dd0219f862b040ebb7037))
* **spider:** standardize metrics naming and add labels ([c40ffba](https://github.com/Leon180/tabelogo-v2/commit/c40ffba5f17e2074069f9fbbbefc74abf1776375))


### Bug Fixes

* **auth:** update tests for new JWT signature ([b8498cf](https://github.com/Leon180/tabelogo-v2/commit/b8498cfe4c3f53911eac4f5372a390f74a8355d0))
* correctly import param ([90e8835](https://github.com/Leon180/tabelogo-v2/commit/90e8835e2eb12b998ca8e7fab5b041dacbc522ea))
* **docker-compose:** unify Redis DB to 0 for session sharing ([92f9e29](https://github.com/Leon180/tabelogo-v2/commit/92f9e295e667dcb2fe2d0782589b42f2a7ba9e5e))
* **frontend:** add authentication to Restaurant Service API client ([f3290d2](https://github.com/Leon180/tabelogo-v2/commit/f3290d2f927a3a76c4fe76245f792161d657bc17))
* **frontend:** add SSE authentication support for Spider Service ([5dfaaba](https://github.com/Leon180/tabelogo-v2/commit/5dfaaba1927b23cf0fadb77ca063a535281d3c41))
* **frontend:** enhance Restaurant Service error handling for auth ([0dbd10c](https://github.com/Leon180/tabelogo-v2/commit/0dbd10c46358dba40ebe170df3ac7c945f1f0108))
* **middleware:** handle JSON string session format from Redis ([553ad08](https://github.com/Leon180/tabelogo-v2/commit/553ad08329c3bca215cec586db163022d5c6428b))
* **restaurant:** relax PATCH permission to authenticated users ([7f4c018](https://github.com/Leon180/tabelogo-v2/commit/7f4c0183a3a89f3eb50349a0f53dec2395cbb310))
* **spider:** add Redis client provider to infrastructure module ([1175184](https://github.com/Leon180/tabelogo-v2/commit/11751847bc29e17fa4328962381f55ed076b377b))
* **spider:** correct pkgconfig import alias ([809261a](https://github.com/Leon180/tabelogo-v2/commit/809261afa77e725576ac8dd9ea6419568765f22f))
* **spider:** correct scraper type in JobProcessor provider ([a0ed055](https://github.com/Leon180/tabelogo-v2/commit/a0ed0556f07cd35e33ca6707a726524d44f1a094))
* **spider:** inject workerCount from config into JobProcessor ([467a9ce](https://github.com/Leon180/tabelogo-v2/commit/467a9cefa0e99c99225191d9c0f5d984c6ae6f9c))
* **spider:** remove duplicate wg.Done() in worker function ([fabab23](https://github.com/Leon180/tabelogo-v2/commit/fabab239362b6c708e6abce9fef2efe08765bf12))
* **spider:** resolve WaitGroup negative counter panic ([6d355b5](https://github.com/Leon180/tabelogo-v2/commit/6d355b53375aaaa1d98ddf1e710d4cba8bca6824))
* **spider:** use background context for job processor workers ([2056635](https://github.com/Leon180/tabelogo-v2/commit/20566355c8746782190911418761b5db7f8a9704))
* **spider:** use GetRedisAddr() method for Redis client ([33a4909](https://github.com/Leon180/tabelogo-v2/commit/33a490904e973779a11f29228a46889f9a5c763e))
* **spider:** use newJobProcessor provider in Module ([5d31754](https://github.com/Leon180/tabelogo-v2/commit/5d317548eb970716c7b337a79ca0db817277a2dc))

## [0.8.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2-v0.7.0...tabelogo-v2-v0.8.0) (2025-12-14)


### Features

* Add AddressComponent to protobuf definition ([cd8554e](https://github.com/Leon180/tabelogo-v2/commit/cd8554e1641e25650d905e11a2bbb365d9af0249))
* Add area column migration for restaurants ([a02f4da](https://github.com/Leon180/tabelogo-v2/commit/a02f4da17e18d16e1ee2b7fbaec8767d1de68140))
* Add area field to restaurant schema ([63223be](https://github.com/Leon180/tabelogo-v2/commit/63223be6147f1b459f6a11139d24da1fdea43148))
* Add area mapper for Tabelog area code conversion ([9b484ae](https://github.com/Leon180/tabelogo-v2/commit/9b484ae3c7478903b5c7df88410aa1a53241cf80))
* Add area parameter to Restaurant constructors and DTO ([ff589f8](https://github.com/Leon180/tabelogo-v2/commit/ff589f86ce0ec7498efaf90c4a3bfdc13f7c2b9d))
* Add CORS support to Spider Service ([3966052](https://github.com/Leon180/tabelogo-v2/commit/3966052e59b916631018a44f487b1a9b26701b81))
* Add frontend API functions for updating restaurant ([8936be1](https://github.com/Leon180/tabelogo-v2/commit/8936be1c7a2b1759826b1be61b660c8f6effa2ce))
* Add gomock integration for auto-generated mocks ([35560e2](https://github.com/Leon180/tabelogo-v2/commit/35560e215d1c0ee252c188392ecbe84ca70fd5f0))
* Add Grafana dashboard and load testing infrastructure ([32425a0](https://github.com/Leon180/tabelogo-v2/commit/32425a0fd5a92aa68a9420b7404e8ea6fb5ef280))
* Add Japanese name support to Restaurant model ([64694ce](https://github.com/Leon180/tabelogo-v2/commit/64694cec4fa9312567efde43ac8616ba07b8e674))
* Add Japanese name support to Spider Service ([b923cb2](https://github.com/Leon180/tabelogo-v2/commit/b923cb29c4e2f698dec77dafa773ff2e2d69dad6))
* Add name_ja to Restaurant DTOs ([442433e](https://github.com/Leon180/tabelogo-v2/commit/442433e30a479aa3d66ffcfa3dc55623c5e9ece2))
* Add PATCH /restaurants/:id endpoint for updating Japanese name ([ba913b5](https://github.com/Leon180/tabelogo-v2/commit/ba913b589b9825a1980a7f2de4d42766b82e8083))
* Add Prometheus metrics for Restaurant Service cache performance ([66f50b7](https://github.com/Leon180/tabelogo-v2/commit/66f50b7c0e1dd6eb96333397a8a993b1e936dfad))
* Add Spider Service Docker integration ([802b49e](https://github.com/Leon180/tabelogo-v2/commit/802b49e05ead5871195e079942d17a94a12f8351))
* Add Spider Service gRPC API ([7103bd1](https://github.com/Leon180/tabelogo-v2/commit/7103bd14b8d368bed9e883a0ca06640cf57cbddd))
* Add Spider Service MVP with DDD architecture ([d284273](https://github.com/Leon180/tabelogo-v2/commit/d284273b07702dd6970126132657bfa8c3abd6e0))
* Add Spider Service to Docker Compose ([ddd8e13](https://github.com/Leon180/tabelogo-v2/commit/ddd8e13b6203769cc7f564e58daad8cc48d6c6eb))
* Add Tabelog button to PlaceDetailModal ([62612ec](https://github.com/Leon180/tabelogo-v2/commit/62612ec9acfe43313eb7ea7313a63a74b9b650f9))
* Complete data mapping from Map Service to Restaurant Service ([f53d09a](https://github.com/Leon180/tabelogo-v2/commit/f53d09a4eff60a166d0a3f1f6db27abb5ca4f512))
* Complete frontend integration prep for Restaurant Service ([94d569f](https://github.com/Leon180/tabelogo-v2/commit/94d569f685d87c13fb0e7717d014413cb38b7910))
* Complete Map Service Phase 1 + Start Phase 2 Integration ([7739085](https://github.com/Leon180/tabelogo-v2/commit/773908597ee87969db5889f983e14fff3abcd90d))
* Complete Phase 2 Restaurant-Map Integration (100%) ([2d03b18](https://github.com/Leon180/tabelogo-v2/commit/2d03b186fe70a773055ad33c5c9ee83d1be723bd))
* Complete Spider Service deployment ([9721fff](https://github.com/Leon180/tabelogo-v2/commit/9721fff54b3835da143f42d539f94dde295d6ab1))
* Complete Spider Service frontend integration ([3510b31](https://github.com/Leon180/tabelogo-v2/commit/3510b31f7a22cef5a1bf527555dd5fdc634a7873))
* Complete SSE implementation with Redis caching ([1d8964b](https://github.com/Leon180/tabelogo-v2/commit/1d8964b57a75593b1b0392f2c0b7ea991cc23ef1))
* Extract area from Google Maps addressComponents ([a016cb5](https://github.com/Leon180/tabelogo-v2/commit/a016cb54f2f0356935e37adacea4eeb529c7c8c2))
* **frontend:** add real-time scraping status display with SSE ([b25ffcd](https://github.com/Leon180/tabelogo-v2/commit/b25ffcdd9c5cb5cbab077a4b0e14a496c0f49206))
* Implement SSE with result caching for Spider Service ([617891c](https://github.com/Leon180/tabelogo-v2/commit/617891c227b8e54c72575059d0d45de12ee11f69))
* Phase 2 Restaurant-Map Integration (80% Complete) ([2547cac](https://github.com/Leon180/tabelogo-v2/commit/2547cac0c1e1480c609b3843fbb91cc0dd9e51ce))
* **restaurant:** add Docker support with gRPC, Swagger UI, and full containerization ([f1100f9](https://github.com/Leon180/tabelogo-v2/commit/f1100f97722d10c4178ba1f61873ca9e08e12284))
* Set up Grafana and Prometheus monitoring ([e629b22](https://github.com/Leon180/tabelogo-v2/commit/e629b222ec7d7e802f714544c316a220b40a956a))
* **spider:** add panic recovery to all goroutines ([84c2d9b](https://github.com/Leon180/tabelogo-v2/commit/84c2d9bf8e2f3ac0ccb4bb8f2e2f4d07004a79d9))
* **spider:** add Prometheus metrics for observability ([29deed8](https://github.com/Leon180/tabelogo-v2/commit/29deed8bf678e1f75591068321d5fc4edbb46e2a))
* **spider:** complete metrics integration for scraper and circuit breaker ([eb62bd2](https://github.com/Leon180/tabelogo-v2/commit/eb62bd2e01f55a9110582d100d5c6da63cd7d2b1))
* **spider:** implement circuit breaker for scraper ([ae895d8](https://github.com/Leon180/tabelogo-v2/commit/ae895d85bfadfea0999ce7269ef5ae998836df78))
* **spider:** implement Phase 1 async job processing with worker pool ([cd96aa7](https://github.com/Leon180/tabelogo-v2/commit/cd96aa727ce2200c0f26bf1a20d2ab8a0d259f01))
* **spider:** implement retry logic with error classification ([eed6ef5](https://github.com/Leon180/tabelogo-v2/commit/eed6ef59a255c3f151c3c8e0b64b5ba77a1a3ce0))
* **spider:** integrate metrics into job processor ([8c5af23](https://github.com/Leon180/tabelogo-v2/commit/8c5af231dbaa045f00f02958eec605da4a23d445))
* **spider:** integrate Phase 1 components into main application (Phase 2) ([414c15a](https://github.com/Leon180/tabelogo-v2/commit/414c15a1f512032b988899b74b4297a206c73038))
* **web:** update spider-service API client for SSE streaming ([0296f98](https://github.com/Leon180/tabelogo-v2/commit/0296f98c7d6817fe33d23659ee1e39a21ec22b05))


### Bug Fixes

* Actually add null check for tabelogResults (previous commit didn't apply) ([2684e7b](https://github.com/Leon180/tabelogo-v2/commit/2684e7b0fcdac288b1ba10de7e65408e84dd390d))
* Add addressComponents extraction in convertToProtoPlace ([cfb0234](https://github.com/Leon180/tabelogo-v2/commit/cfb02345f4f6d34f8d13b5483a1c1ae976d1e5dc))
* Add addressComponents to API mask in Restaurant Service ([397adf1](https://github.com/Leon180/tabelogo-v2/commit/397adf1966f395037dcc0cdf78b6e3d285650cc7))
* Add addressComponents to Google Maps API field mask ([1883cbf](https://github.com/Leon180/tabelogo-v2/commit/1883cbf39a1e031ec7924eb252151e693328be70))
* Add addressComponents to Place type definition ([746de74](https://github.com/Leon180/tabelogo-v2/commit/746de740b24d6b372618a8c4b163af53e9b4a4ec))
* Add go mod tidy to all service Dockerfiles ([fd1c28c](https://github.com/Leon180/tabelogo-v2/commit/fd1c28ccef64c5d5b6864b045630bcfc61bd1343))
* Add missing update methods to Restaurant model ([6ea1dfa](https://github.com/Leon180/tabelogo-v2/commit/6ea1dfad7650b2221f29f4fac138ee65506820e1))
* Add NameJa field to application DTO and service layer ([d15c8d8](https://github.com/Leon180/tabelogo-v2/commit/d15c8d8fa2f20e50c4be865d5cd0a6bab9d7f92a))
* Add NoRoute handler for CORS OPTIONS requests ([9f4f203](https://github.com/Leon180/tabelogo-v2/commit/9f4f20323425ae9312670fad5bdc069654a28706))
* Add null check for tabelogResults ([17ad1ab](https://github.com/Leon180/tabelogo-v2/commit/17ad1abfee47beb5f4204db53591ce47943afd87))
* Add null safety checks for Tabelog restaurant types and photos ([9b46fe1](https://github.com/Leon180/tabelogo-v2/commit/9b46fe18039a46eb9053a684e1f96cc194b5a323))
* Add optional chaining for tabelogRestaurants to prevent undefined error ([c8a1060](https://github.com/Leon180/tabelogo-v2/commit/c8a10602549039b0458582c0b4422e62b9a7fa1d))
* Change healthcheck from HEAD to GET for all services ([3e60a10](https://github.com/Leon180/tabelogo-v2/commit/3e60a10e03a1d5efafd98a9590476962e7db8b21))
* Change Spider Service env var from HTTP_PORT to SERVER_PORT ([2d1031c](https://github.com/Leon180/tabelogo-v2/commit/2d1031cfbd981559a82addb80162f10c91bc328d))
* Complete Spider server update for Japanese name support ([cc4d58a](https://github.com/Leon180/tabelogo-v2/commit/cc4d58ae41dcaa38327f60fae4c371d59d42465d))
* Correct function calls and add area to restaurant update ([5d81d2c](https://github.com/Leon180/tabelogo-v2/commit/5d81d2c1cc9b9924aed90498a1928c9d3036ba31))
* Correct MapServiceClient mock generation ([6f92ea1](https://github.com/Leon180/tabelogo-v2/commit/6f92ea1a21740ce2a9b923f6b1840437bd608922))
* Correct method calls in Spider server ([e855b51](https://github.com/Leon180/tabelogo-v2/commit/e855b51bc136dff10d02f6861ffeefbe98125e71))
* Correct Spider Service API endpoint from /search to /scrape ([815912a](https://github.com/Leon180/tabelogo-v2/commit/815912ae08c2a128671fa294de935ad60e07ed75))
* Correct SSE status comparison to uppercase (COMPLETED/FAILED) ([d5aa143](https://github.com/Leon180/tabelogo-v2/commit/d5aa143cda18a933182b46fc01218e382474f4f4))
* Correct Tabelog results display with proper variable name ([79f33c1](https://github.com/Leon180/tabelogo-v2/commit/79f33c11b13bf46b7edfa7f77b2bfd7208da2724))
* Correct Tabelog URL construction for search ([7e3ff1c](https://github.com/Leon180/tabelogo-v2/commit/7e3ff1cd1a8bac7d0976ff5cb6e64eb4082e90f0))
* Fetch addressComponents when getting Japanese name ([35a0a3d](https://github.com/Leon180/tabelogo-v2/commit/35a0a3d2d5867eacde762a09e1e129a2adc68f69))
* Force English language code for addressComponents in Map Service ([22fa841](https://github.com/Leon180/tabelogo-v2/commit/22fa841ca7e8048e789d772facd4aacdd17d8638))
* **frontend:** add default case to ScrapingStatus config ([4b289a2](https://github.com/Leon180/tabelogo-v2/commit/4b289a20fe941a449f81053b79ae9a003940322a))
* **frontend:** normalize SSE status to lowercase ([5e70e9e](https://github.com/Leon180/tabelogo-v2/commit/5e70e9e9d44c4fddcdc8e9c7e3d2c54710843c60))
* **frontend:** remove leftover setScrapingProgress call ([efeaa0f](https://github.com/Leon180/tabelogo-v2/commit/efeaa0f570d4bf94903807380720f2157904555d))
* Make TabelogRestaurant fields public for JSON serialization ([7c2a4f7](https://github.com/Leon180/tabelogo-v2/commit/7c2a4f71da2fa29073e32e404665739b4088aad2))
* Map NameJa field in HTTP handler ([e7183f9](https://github.com/Leon180/tabelogo-v2/commit/e7183f96de8a93bf20c380bc064594759112f1d3))
* Match auth-service CORS pattern with explicit OPTIONS routes ([e6d6188](https://github.com/Leon180/tabelogo-v2/commit/e6d61885acac3496838634d44681f65a462951c2))
* Prevent SSE subscription when cache hit occurs ([a961c57](https://github.com/Leon180/tabelogo-v2/commit/a961c57ac944102553a2dbd16fcae9530d00a480))
* Prioritize English language code in area extraction ([84fabf3](https://github.com/Leon180/tabelogo-v2/commit/84fabf37fd33b89ca9b977eca8bb552988192dd3))
* Rename cache response field from 'results' to 'restaurants' ([7a6abdf](https://github.com/Leon180/tabelogo-v2/commit/7a6abdf1c291802c07036ce509d97d71f9033a6e))
* Resolve Go version mismatch for mock generation ([6bce2d2](https://github.com/Leon180/tabelogo-v2/commit/6bce2d23fcad49d2878c064f42f57fa1642c8a85))
* Restore Go 1.24 and upgrade mockgen for compatibility ([b868dbd](https://github.com/Leon180/tabelogo-v2/commit/b868dbdbdd4c96f530cdeab5d8ea1632b669a731))
* Restore useQuickSearch function that was accidentally broken ([3a738a2](https://github.com/Leon180/tabelogo-v2/commit/3a738a298deef1bf93c11f8a01cae4a74128e166))
* Separate MapServiceClient interface to fix mock generation ([2cd7359](https://github.com/Leon180/tabelogo-v2/commit/2cd7359c0f02edd054853ee0836983c1b02e052c))
* Spider Service healthcheck and port configuration ([e3f4548](https://github.com/Leon180/tabelogo-v2/commit/e3f45485a6bf9f25364fb08166d533d5a2c825df))
* **spider:** complete circuit breaker test fixes ([14c2c9e](https://github.com/Leon180/tabelogo-v2/commit/14c2c9ee22416321f69bf773893ec6c812c132b4))
* **spider:** correct parameter order in ScrapeRestaurants call ([8904f22](https://github.com/Leon180/tabelogo-v2/commit/8904f22398e4301dbfc67f3846d267e806e80b03))
* **spider:** correct syntax error in infrastructure module ([b970025](https://github.com/Leon180/tabelogo-v2/commit/b970025f6bcc31ad41892ad213c638ad58dddfa8))
* **spider:** correct syntax error in module.go ([4f53558](https://github.com/Leon180/tabelogo-v2/commit/4f5355800ccedda79ddbeeee65949042849a02a4))
* **spider:** implement DTO pattern for Redis cache serialization ([beabd6a](https://github.com/Leon180/tabelogo-v2/commit/beabd6a617cab0a3f3b1e5e7a23d76a9f624ef82))
* **spider:** include results data in SSE updates ([eb6018f](https://github.com/Leon180/tabelogo-v2/commit/eb6018f447adfc77383e3e82dcffb1657e39518d))
* **spider:** integrate JobProcessor with use case for proper async processing ([fb27d99](https://github.com/Leon180/tabelogo-v2/commit/fb27d99fd2a06e228aa67041687ee4776add271e))
* **spider:** remove duplicate /metrics endpoint registration ([c5428b1](https://github.com/Leon180/tabelogo-v2/commit/c5428b110345ac5fac30c6317cbc43d87a581570))
* **spider:** remove NewRedis and fix type references ([a1b1454](https://github.com/Leon180/tabelogo-v2/commit/a1b1454dcbb0f730fa2faf0d967e72f97430750e))
* **spider:** remove non-existent NewRateLimiter from module ([c11daf1](https://github.com/Leon180/tabelogo-v2/commit/c11daf1d6459b17433ca6e4ea1451af6947bf9d7))
* **spider:** resolve type reference errors in infrastructure module ([76b66a5](https://github.com/Leon180/tabelogo-v2/commit/76b66a5d24d41b9419ce135f10fb67a0b445cecc))
* **spider:** update circuit breaker tests for metrics parameter ([bc6a923](https://github.com/Leon180/tabelogo-v2/commit/bc6a92305aa00db5bde00c8301286197d5f5362c))
* **spider:** use background context for worker pool lifetime ([3184c4b](https://github.com/Leon180/tabelogo-v2/commit/3184c4b2d184f69b372a973b44dba27e9db4a91c))
* **spider:** use correct types in infrastructure module ([d487e0d](https://github.com/Leon180/tabelogo-v2/commit/d487e0d12d15f967eae4ae8b457eb654ad4d5c83))
* **spider:** use NewScraperConfig with builder pattern ([5b33ca9](https://github.com/Leon180/tabelogo-v2/commit/5b33ca9b937348a385c7fc6d282f7669c8e22a40))
* **spider:** use RecordScrapeError instead of non-existent RecordJobError ([118d529](https://github.com/Leon180/tabelogo-v2/commit/118d529caf1a01b63f52494f5744be166ac5c7d0))
* Update all unit tests with new RestaurantService constructor ([76bb92d](https://github.com/Leon180/tabelogo-v2/commit/76bb92dc81465d7be60d921efa1c132fac966e85))
* Update gRPC UpdateRestaurant to use new DTO structure ([8eb54c8](https://github.com/Leon180/tabelogo-v2/commit/8eb54c847ed05e267eacc9c1d5b005c2b3a29b96))
* Update Makefile build target to use docker-compose ([fae418d](https://github.com/Leon180/tabelogo-v2/commit/fae418dc9819e86b694b0ff76b247bb664171d9a))
* Update proto and gRPC server for UpdateRestaurant ([91fadfe](https://github.com/Leon180/tabelogo-v2/commit/91fadfeb2f025d48a5842fd3d75713644e5858fc))
* Update Spider Service dependencies ([54d6f01](https://github.com/Leon180/tabelogo-v2/commit/54d6f01a5ace8990e5ae6483e6214e5089fd3621))
* Update Spider Service Dockerfile to match auth-service pattern ([e2b2f49](https://github.com/Leon180/tabelogo-v2/commit/e2b2f49013061511aef1fe0c309b66df32b034f0))
* Update TypeScript status types to match backend enum ([49e9333](https://github.com/Leon180/tabelogo-v2/commit/49e9333f30253d4ee4084f1726c827e48c1a696a))
* Use English for quick search, Japanese only for Tabelog ([992375f](https://github.com/Leon180/tabelogo-v2/commit/992375f78e554384312e095a2c161bcd09006455))
* Use JobStatus enum constants for cache and SSE logic ([f22fb6d](https://github.com/Leon180/tabelogo-v2/commit/f22fb6da8a14c6a5583512b3d5e07caa9474814f))


### Performance Improvements

* Optimize Docker builds with BuildKit cache ([cd9b861](https://github.com/Leon180/tabelogo-v2/commit/cd9b8618d52dee4ab0cdd8bbe50c323f54333e4d))

## [0.7.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2-v0.6.0...tabelogo-v2-v0.7.0) (2025-12-03)


### Features

* **restaurant:** complete Restaurant Service implementation with 98% test coverage ([e831f76](https://github.com/Leon180/tabelogo-v2/commit/e831f7618112e7bd8101504b9426e50adbad0f26))

## [0.6.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2-v0.5.0...tabelogo-v2-v0.6.0) (2025-12-02)


### Features

* **frontend:** add Quick Search modal and POI click functionality ([0a8d038](https://github.com/Leon180/tabelogo-v2/commit/0a8d0385f2c13bccdb20ef6680ec3d8fc923ccfa))
* **frontend:** add restaurant list view with resizable sidebar and collapsible search ([390cc4b](https://github.com/Leon180/tabelogo-v2/commit/390cc4b3fc4f73ae825d4c0795f7f3ea56997166))
* **frontend:** integrate Map Service with React Query and TypeScript ([65e1a4c](https://github.com/Leon180/tabelogo-v2/commit/65e1a4c7f4f144c11f646dc53360f392e7c03c08))
* init plan ([e3eb7d2](https://github.com/Leon180/tabelogo-v2/commit/e3eb7d21d6632c3dc46c64bf336a497d7357d6f2))
* **map-service:** implement complete Map Service with Phase 1-4 ([c4670ee](https://github.com/Leon180/tabelogo-v2/commit/c4670ee0675de81ae24e8a0bfec132ff5489079a))


### Bug Fixes

* add null checks for place.location in GoogleMap component ([9f09c3d](https://github.com/Leon180/tabelogo-v2/commit/9f09c3d2af0a9e57646c29c4318d55b82c69b978))

## [0.5.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2-v0.4.1...tabelogo-v2-v0.5.0) (2025-11-25)


### Features

* implement user authentication with login UI and CORS support ([59db989](https://github.com/Leon180/tabelogo-v2/commit/59db989597c75d5406652cf0284a3aea88a7ef26))

## [0.4.1](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2-v0.4.0...tabelogo-v2-v0.4.1) (2025-11-24)


### Bug Fixes

* **auth-service:** resolve Swagger UI access issues in local and Docker environments ([b13ff12](https://github.com/Leon180/tabelogo-v2/commit/b13ff12f43ad27e4dd12af39eb4e46050fbec299))

## [0.4.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2-v0.3.0...tabelogo-v2-v0.4.0) (2025-11-24)


### Features

* **frontend:** implement Next.js frontend with Google Maps integration ([5be5847](https://github.com/Leon180/tabelogo-v2/commit/5be5847a6a9609f66d53be173cdabdaeebff4bbc))

## [0.3.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2-v0.2.0...tabelogo-v2-v0.3.0) (2025-11-24)


### Features

* **auth:** finalize docker and makefile configuration ([d84778a](https://github.com/Leon180/tabelogo-v2/commit/d84778a50322b99f1e1ab5bcfe7488d3eb2515ae))
* **auth:** implement domain, infra, app, and grpc layers ([ce6abd7](https://github.com/Leon180/tabelogo-v2/commit/ce6abd7ae9dd67fab16e19c0b96cfd677e7455a6))
* update architecture.md ([6fe0376](https://github.com/Leon180/tabelogo-v2/commit/6fe0376bdf9b33c04242fd9e0c332aa1feef4152))

## [0.2.0](https://github.com/Leon180/tabelogo-v2/compare/tabelogo-v2-v0.1.0...tabelogo-v2-v0.2.0) (2025-11-22)


### Features

* allow prefix env key ([9a8f6dd](https://github.com/Leon180/tabelogo-v2/commit/9a8f6dd07fb9538f9476badf07d7e3df7572ca68))
* complete Phase 1 - migrations, shared packages, and middleware ([6e4933c](https://github.com/Leon180/tabelogo-v2/commit/6e4933c132eb48d7f431926f4d34ad8084c907e5))
* init ([532a5bf](https://github.com/Leon180/tabelogo-v2/commit/532a5bfdfeefdd3d2404410e99d84061e8a8a5e2))
* init migrations pkg and test case ([7b624f7](https://github.com/Leon180/tabelogo-v2/commit/7b624f76113d03ce02aadb12395f38a435f09be8))
* setup release-please for monorepo versioning ([fbcc243](https://github.com/Leon180/tabelogo-v2/commit/fbcc243e0b5497b4a31abde4c22311ba38025653))


### Bug Fixes

* modify mod name, and apply fx to pkg ([960e177](https://github.com/Leon180/tabelogo-v2/commit/960e17712c30c9b162551f94aed0c38b6942f1af))
* standardize Go version to 1.23 and update gitignore for workspace files ([05f10d9](https://github.com/Leon180/tabelogo-v2/commit/05f10d9f5985e1de842886b66bef0c6a5b5aabc6))
