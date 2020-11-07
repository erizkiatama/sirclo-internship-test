import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

import static org.junit.jupiter.api.Assertions.*;

class CartTest {

    Cart keranjang;

    @BeforeEach
    void setUp() {
        keranjang = new Cart();
    }

    @Test
    void testTambahProdukWhenTheProductNotInTheCart_ShouldEnterTheProductInToTheCart() {
        keranjang.tambahProduk("Air Mineral", 3);

        assertEquals(1, keranjang.getKeranjang().size());
        assertEquals(3, keranjang.getKeranjang().get("Air Mineral"));
    }

    @Test
    void testTambahProdukWhenTheProductAlreadyInTheCart_ShouldAddQuantityOfThatProduct() {
        keranjang.tambahProduk("Air Mineral", 3);
        assertEquals(1, keranjang.getKeranjang().size());
        assertEquals(3, keranjang.getKeranjang().get("Air Mineral"));


        keranjang.tambahProduk("Air Mineral", 7);
        assertEquals(1, keranjang.getKeranjang().size());
        assertEquals(10, keranjang.getKeranjang().get("Air Mineral"));
    }

    @Test
    void testTambahProdukWithMoreThanOneProduct_ShouldEnterAllOfTheProductsInToTheCart() {
        keranjang.tambahProduk("Air Mineral", 3);
        keranjang.tambahProduk("Snack Micin", 5);

        assertEquals(2, keranjang.getKeranjang().size());
        assertEquals(3, keranjang.getKeranjang().get("Air Mineral"));
        assertEquals(5, keranjang.getKeranjang().get("Snack Micin"));
    }

    @Test
    void testHapusProdukWhenNoProductInTheCart_ShouldDoNothing() {
        assertEquals(0, keranjang.getKeranjang().size());

        keranjang.hapusProduk("Air Mineral");
        assertEquals(0, keranjang.getKeranjang().size());
    }

    @Test
    void testHapusProdukWhenKodeProdukNoMatchesWithAnyProduct_ShouldDoNothing() {
        keranjang.tambahProduk("Air Mineral", 3);
        assertEquals(1, keranjang.getKeranjang().size());

        keranjang.hapusProduk("Bukan Air Mineral");
        assertEquals(1, keranjang.getKeranjang().size());
    }

    @Test
    void testHapusProdukWhenKodeProdukMatchesAnyProduct_ShouldDeleteThatProduct() {
        keranjang.tambahProduk("Air Mineral", 3);
        assertEquals(1, keranjang.getKeranjang().size());

        keranjang.hapusProduk("Air Mineral");
        assertEquals(0, keranjang.getKeranjang().size());
        assertNull(keranjang.getKeranjang().get("Air Mineral"));
    }

    @Test
    void testTampilkanCartWhenNoProductInCart_ShouldPrintCartEmptyMessage() throws Exception {
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        System.setOut(new PrintStream(out));

        keranjang.tampilkanCart();

        String expectedMessage = "Keranjang belanja anda kosong\n";

        assertEquals(expectedMessage, out.toString());
    }

    @Test
    void testTampilkanCartWhenThereAreProductsInCart_ShouldPrintTheProductWithGivenFormat() {
        keranjang.tambahProduk("Air Mineral", 3);
        keranjang.tambahProduk("Snack Micin", 5);

        ByteArrayOutputStream out = new ByteArrayOutputStream();
        System.setOut(new PrintStream(out));

        keranjang.tampilkanCart();

        String expectedMessage = "Air Mineral (3)\nSnack Micin (5)\n";

        assertEquals(expectedMessage, out.toString());
    }
}