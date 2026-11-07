import XCTest

/// README screenshots of the iPad app, in English and German: scan result
/// (a medicine bottle from the photo library, put there with
/// `simctl addmedia`), history, create (with a QR code), about.
///
/// The PNGs go to $SHOTS_DIR (passed as TEST_RUNNER_SHOTS_DIR by xcodebuild);
/// the simulator writes straight into the Mac's file system.
final class Screenshots: XCTestCase {
    let app = XCUIApplication(bundleIdentifier: "eu.mlcgo.mlcbarcode")

    var shotsDir: URL {
        URL(fileURLWithPath: ProcessInfo.processInfo.environment["SHOTS_DIR"] ?? NSTemporaryDirectory())
    }

    struct Lang {
        let code, scan, history, create, about, choose, make, scanAgain, language, expect: String
    }
    let langs = [
        Lang(code: "en", scan: "Scan", history: "History", create: "Create", about: "About",
             choose: "Choose photo", make: "Create QR code", scanAgain: "Scan again",
             language: "English", expect: "12345678"),
        Lang(code: "de", scan: "Scannen", history: "Verlauf", create: "Erzeugen", about: "Über",
             choose: "Foto wählen", make: "QR-Code erzeugen", scanAgain: "Weiter scannen",
             language: "Deutsch", expect: "12345678"),
    ]

    override func setUp() {
        continueAfterFailure = false
    }

    /// The first element whose label is `text` or ends with it — the tab
    /// buttons carry an icon glyph in front of their text.
    func element(_ text: String, _ type: XCUIElement.ElementType = .any) -> XCUIElement {
        let pred = NSPredicate(format: "label == %@ OR label ENDSWITH %@", text, " " + text)
        return app.descendants(matching: type).matching(pred).firstMatch
    }

    func tap(_ text: String, _ type: XCUIElement.ElementType = .any, timeout: TimeInterval = 15) {
        let e = element(text, type)
        XCTAssertTrue(e.waitForExistence(timeout: timeout), "not on screen: \(text)")
        e.tap()
        sleep(1)
    }

    func shot(_ name: String) {
        let png = XCUIScreen.main.screenshot().pngRepresentation
        try? FileManager.default.createDirectory(at: shotsDir, withIntermediateDirectories: true)
        XCTAssertNoThrow(try png.write(to: shotsDir.appendingPathComponent("\(name).png")))
    }

    /// <input type="file"> on iOS first offers a menu (photo library,
    /// camera, files), then the photo picker; pick the newest photo.
    func pickPhoto() {
        let library = app.descendants(matching: .any).matching(
            NSPredicate(format: "label IN %@", ["Photo Library", "Fotomediathek", "Fotos"])).firstMatch
        if library.waitForExistence(timeout: 5) { library.tap() }
        let photo = app.images.matching(NSPredicate(format: "label BEGINSWITH 'Photo' OR label BEGINSWITH 'Foto'")).firstMatch
        XCTAssertTrue(photo.waitForExistence(timeout: 15), "no photo in the picker")
        photo.tap()
        // The picker then shows the selection and waits for a confirmation.
        let done = app.buttons.matching(NSPredicate(format: "label IN %@", ["Add", "Hinzufügen", "Done", "Fertig"])).firstMatch
        if done.waitForExistence(timeout: 5) { done.tap() }
    }

    func testReadmeScreenshots() {
        app.launch()
        for l in langs {
            // Switch language on the About tab, in whichever language is on.
            let about = [element("About", .button), element("Über", .button)]
            XCTAssertTrue(about[0].waitForExistence(timeout: 30) || about[1].exists, "no About tab")
            (about[0].exists ? about[0] : about[1]).tap()
            sleep(1)
            tap(l.language)
            tap(l.scan, .button)
            // A result from the previous language still covers the button.
            let again = element(l.scanAgain)
            if again.waitForExistence(timeout: 3) { again.tap(); sleep(1) }
            tap(l.choose)
            pickPhoto()
            let result = app.descendants(matching: .any).matching(NSPredicate(format: "label CONTAINS %@", l.expect)).firstMatch
            if !result.waitForExistence(timeout: 20) {
                shot("fail-\(l.code)-scan")
                XCTFail("scan result missing")
            }
            sleep(1)
            shot("ipad-\(l.code)-scan")

            tap(l.history, .button)
            shot("ipad-\(l.code)-history")

            tap(l.create, .button)
            let field = app.textViews.firstMatch
            XCTAssertTrue(field.waitForExistence(timeout: 10))
            field.tap()
            field.typeText("https://mlcgo.eu/products/mlc-barcode/")
            tap(l.make)
            sleep(2)
            shot("ipad-\(l.code)-create")

            tap(l.about, .button)
            shot("ipad-\(l.code)-about")
        }
    }
}
